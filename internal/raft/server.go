package raft

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type NodeServer struct {
	Addr        string
	IdleTimeout time.Duration
	inShutdown  bool
	listener    net.Listener
	conns []*Conn
}

func (srv *NodeServer) ListenAndServe() error {
	var err error
	addr := srv.Addr

	if addr == "" {
		addr = ":6000"
	}

	fmt.Printf("Start server on %s\n", srv.Addr)

	srv.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer srv.listener.Close()

	for {
		if srv.inShutdown {
			continue
		}
		newConn, err := srv.listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection %v", err)
			continue
		}

		conn := &Conn{
			Conn:        newConn,
			IdleTimeout: srv.IdleTimeout,
		}

		_ = conn.SetDeadline(time.Now().Add(conn.IdleTimeout))

		srv.conns = append(srv.conns, conn)

		fmt.Printf("accepted connection from %v", conn.RemoteAddr())
		go srv.handle(conn)
	}
}

func (srv *NodeServer) handle(conn net.Conn) {
	defer func() {
		fmt.Printf("closing connection from %v\n", conn.RemoteAddr())
		_ = conn.Close()
	}()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(r)
	for {
		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				fmt.Printf("%v(%v)\n", err, conn.RemoteAddr())
				return
			}
			break
		}
		_, _ = w.WriteString(strings.ToUpper(scanner.Text()) + "\n")
		_ = w.Flush()
	}
}

func (srv *NodeServer) Shutdown() {
	srv.inShutdown = true

	fmt.Println("shutting down...")

	_ = srv.listener.Close()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("waiting on %v connections", len(srv.conns))
		}
		if len(srv.conns) == 0 {
			return
		}
	}
}

type Conn struct {
	net.Conn
	IdleTimeout   time.Duration
	MaxReadBuffer int64
}

func (c *Conn) Write(p []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Write(p)
}

func (c *Conn) Read(b []byte) (int, error) {
	c.updateDeadline()
	r := io.LimitReader(c.Conn, c.MaxReadBuffer)
	return r.Read(b)
}

func (c *Conn) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	_ = c.Conn.SetDeadline(idleDeadline)
}
