package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func selectMenu(servers []*server) *server {
	var options []string
	for _, server := range servers {
		options = append(options, server.String())
	}

	var selected string
	if err := survey.AskOne(&survey.Select{
		Message: "Select a server",
		Options: options,
	}, &selected); err != nil {
		return nil
	}

	for index, option := range options {
		if selected == option {
			return servers[index]
		}
	}

	return nil
}

func selectArg(servers []*server) *server {
	query := strings.ToLower(strings.Join(os.Args[1:], " "))
	for _, server := range servers {
		for _, alias := range server.Aliases {
			if strings.ToLower(alias) == query {
				return server
			}
		}

		if strings.ToLower(server.Name) == query || strings.ToLower(server.Hostname) == query {
			return server
		}
	}

	return nil
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		println("Failed to get home directory:", err.Error())
		return
	}

	servers, err := loadServers(filepath.Join(home, ".smm.yml"))
	if err != nil {
		println("Failed to load servers:", err.Error())
		return
	}

	if len(servers) == 0 {
		println("Sorry, there are no servers yet =(")
		return
	}

	srv := selectArg(servers)
	if srv == nil {
		srv = selectMenu(servers)
	}

	if srv == nil {
		return
	}

	cmdName := "ssh"
	var cmdArgs []string

	if srv.Password != "" {
		cmdName = "sshpass"
		cmdArgs = append(cmdArgs, "-p", srv.Password, "ssh")
	}

	cmdArgs = append(cmdArgs, srv.Username+"@"+srv.Hostname)
	cmdArgs = append(cmdArgs, "-p", strconv.Itoa(srv.port))
	if srv.KeyFile != "" {
		cmdArgs = append(cmdArgs, "-i", srv.KeyFile)
	}
	if srv.ForceTTY {
		cmdArgs = append(cmdArgs, "-t")
	}
	if srv.Command != "" {
		cmdArgs = append(cmdArgs, srv.Command)
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err = cmd.Start(); err != nil {
		println("Failed to run command:", err.Error())
		return
	}

	cmd.Wait()
}
