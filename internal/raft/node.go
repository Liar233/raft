package raft

type Node struct {
	config *AppConfig	
}

func NewNode(config *AppConfig) *Node {
	return &Node{
		config: config,
	}
}

