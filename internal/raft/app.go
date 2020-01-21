package raft

import "fmt"

var Config *AppConfig

func StartApp(filename string) {
	var err error
	Config, err = parseConfig(filename)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	node := NewNode(Config)

	err = node.Start()

	if err != nil {
		fmt.Println(err.Error())
	}
}
