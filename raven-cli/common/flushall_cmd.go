package common

import (
	"github.com/abiosoft/ishell"
)

func FlushAllCmd() *ishell.Cmd {
	return flushAllCmd
}

var flushAllCmd = &ishell.Cmd{
	Name:     "flushall",
	Aliases:  []string{},
	Help:     "Flush All MessageBoxes [main + processing + dead]",
	LongHelp: "",
	Func: func(c *ishell.Context) {
		rs, _ := GetRavenServer(c)

		result, err := rs.FlushAll()
		if err != nil {
			c.Println(err.Error())
			return
		}
		c.Println("MessageBoxes Flushed:")
		c.Println("-----------------------------")
		for k, v := range result {
			c.Printf("%s => %v\n", k, v)
		}
		c.Println()
		return
	},
}
