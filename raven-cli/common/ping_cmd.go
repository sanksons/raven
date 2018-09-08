package common

import (
	"github.com/abiosoft/ishell"
)

func GetPingCmd() *ishell.Cmd {
	return pingCmd
}

var pingCmd = &ishell.Cmd{
	Name:     "ping",
	Aliases:  []string{"Ping", "health", "hc"},
	Help:     "Ping Server",
	LongHelp: "",
	Func: func(c *ishell.Context) {
		rs, _ := GetRavenServer(c)
		res, err := rs.Ping()
		if err != nil {
			c.Println("Unable to connect to server")
			c.Println(err.Error())
			return
		}
		c.Printf("Server Said: %s\n", res)
		return

	},
}
