package raft

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

const NodeName = "RAFT_NODE"

type NodeConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	Name            string
	ElectionTimeout time.Duration         `yaml:"election-timeout"`
	Nodes           map[string]NodeConfig `yaml:"nodes"`
}

func parseConfig(filename string) (*AppConfig, error) {
	var err error
	var data []byte
	config := &AppConfig{}

	if _, err = os.Stat(filename); err != nil {
		return nil, fmt.Errorf("config file `%s` does not exist", filename)
	}

	data, err = ioutil.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("can't read `%s` file with error: %s", filename, err)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	name := os.Getenv(NodeName)

	if name == "" {
		return nil, fmt.Errorf("can't find node name in variable %s:%s", NodeName, name)
	}

	config.Name = name

	return config, nil
}
