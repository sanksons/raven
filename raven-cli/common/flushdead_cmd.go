package common

import (
	"github.com/abiosoft/ishell"
)

func FlushDeadCmd() *ishell.Cmd {
	return flushDeadCmd
}

var flushDeadCmd = &ishell.Cmd{
	Name:     "flushdead",
	Aliases:  []string{},
	Help:     "Flush DeadBox Contents",
	LongHelp: "",
	Func: func(c *ishell.Context) {
		rs, _ := GetRavenServer(c)

		result, err := rs.FlushDeadQ()
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
