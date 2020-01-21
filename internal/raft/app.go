package raft

import "fmt"

const Hello = "Hello"

var Config *AppConfig

func StartApp(filename string)  {
	var err error
	Config, err = parseConfig(filename)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%+v", Config)
}
