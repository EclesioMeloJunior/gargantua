package config

import (
	"encoding/json"
	"io/ioutil"
)

type (
	Node struct {
		Protocol string `json:"protocol"`
		Basepath string `json:"basepath"`
	}

	Network struct {
		Port      string   `json:"port"`
		Bootnodes []string `json:"bootnodes"`
	}

	Config struct {
		Node    *Node    `json:"node"`
		Network *Network `json:"network"`
	}
)

func FromJson(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	if err = json.Unmarshal(b, c); err != nil {
		return nil, err
	}

	return c, nil
}
