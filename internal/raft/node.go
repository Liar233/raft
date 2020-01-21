package raft

import "time"

type Node struct {
	config *AppConfig
	server *NodeServer
}

func NewNode(config *AppConfig) *Node {
	return &Node{
		config: config,
		server: NewNodeServer(
			"0.0.0.0:6000",
			time.Duration(5)*time.Second,
			10240,
		),
	}
}

func (node *Node) Start() error {
	return node.server.ListenAndServe()
}
