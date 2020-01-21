package main

import (
	"github.com/Liar233/raft/internal/raft"
	"os"
)

func main() {
	raft.StartApp(os.Args[1])
}
