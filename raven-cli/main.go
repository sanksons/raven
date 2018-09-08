package main

import (
	"flag"
	"fmt"

	"github.com/sanksons/raven/raven-cli/common"

	"github.com/abiosoft/ishell"
)

var cmds []*ishell.Cmd

func main() {

	portPtr := flag.String("port", "", "Port on Which Receiver is running")
	hostPtr := flag.String("host", "", "Host on Which Receiver is running")
	flag.Parse()

	var port string = *portPtr
	if port == "" {
		fmt.Println("No port specified")
		return
	}

	var host string = *hostPtr
	if host == "" {
		host = "127.0.0.1"
	}

	sd := common.RavenServer{IP: host, Port: port}

	//Check if we can ping server.
	_, err := sd.Ping()
	if err != nil {
		fmt.Println("Could not connect to server :(")
		fmt.Println(err.Error())
		return
	}

	shell := sd.GetNewShell()

	shell.AddCmd(common.GetPingCmd())
	shell.AddCmd(common.GetStatusCmd())

	shell.AutoHelp(true)
	//shell.ClearScreen()
	// run shell
	shell.Run()
}

func RegisterCommand(cmd *ishell.Cmd) {
	cmds = append(cmds, cmd)
	return
}
