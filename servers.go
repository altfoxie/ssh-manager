package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

type server struct {
	// Info
	Name    string   `yaml:"name"`
	Aliases []string `yaml:"aliases"`

	// Connection
	Hostname string `yaml:"hostname"`
	Port     *int   `yaml:"port"`

	// Auth
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	KeyFile  string `yaml:"key_file"`

	// Prekoli
	ForceTTY bool   `yaml:"force_tty"`
	Command  string `yaml:"command"`

	// lol
	port int
}

func (s server) String() (result string) {
	result += s.Name
	if len(s.Aliases) > 0 {
		result += " [" + strings.Join(s.Aliases, ", ") + "]"
	}
	result += fmt.Sprintf(" (%s@%s:%d)", s.Username, s.Hostname, s.port)
	return
}

func loadServers(path string) ([]*server, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var servers []*server
	if err = yaml.Unmarshal(body, &servers); err != nil {
		return nil, err
	}

	for index, server := range servers {
		if server.Name == "" {
			return nil, fmt.Errorf("server %d: name cannot be empty", index)
		}

		if server.Hostname == "" {
			return nil, fmt.Errorf("server %d: hostname cannot be empty", index)
		}

		if server.Username == "" {
			return nil, fmt.Errorf("server %d: username cannot be empty", index)
		}

		if server.Port != nil {
			server.port = *server.Port
			if server.port < 1 || server.port > 65535 {
				return nil, fmt.Errorf("server %d: port cannot be less than 1 or greater than 65535", index)
			}
		} else {
			server.port = 22
		}
	}

	return servers, nil
}
