package common

import (
	"fmt"

	"github.com/sanksons/gowraps/filesystem"

	"github.com/abiosoft/ishell"
)

func ShowDeadBoxCmd() *ishell.Cmd {
	return showDeadBoxCmd
}

var showDeadBoxCmd = &ishell.Cmd{
	Name:     "showdead",
	Aliases:  []string{"dead"},
	Help:     "Show Dead Messages",
	LongHelp: "",
	Func: func(c *ishell.Context) {
		c.Print("Save to File: ")
		filepath := c.ReadLine()

		if filepath != "" {
			if err := filesystem.DeleteFile(filepath); err != nil {
				c.Println(fmt.Sprintf("Got Error: %s", err.Error()))
				return
			}
		}

		var showInScreen bool
		c.Print("Show results on screen: Y or N ?")
		show := c.ReadLine()
		if show == "Y" || show == "y" {
			showInScreen = true
		}

		rs, _ := GetRavenServer(c)
		msgs, err := rs.ShowDeadBox()
		if err != nil {
			c.Println("Unable to connect to server")
			c.Println(err.Error())
			return
		}

		for _, m := range msgs {
			data := ""
			data += fmt.Sprintln("---------------------------")
			data += fmt.Sprintln(fmt.Sprintf("Id: %s", m.Id))
			data += fmt.Sprintln(fmt.Sprintf("Type: %s", m.Type))
			data += fmt.Sprintln(fmt.Sprintf("Data: %s", m.Data))
			data += fmt.Sprintln(fmt.Sprintf("ShardKey: %s", m.ShardKey))
			if showInScreen {
				c.Println(data)
			}
			if filepath != "" {
				filesystem.AppendToFile(filepath, []byte(data))
			}
		}
		if filepath != "" {
			c.Println("===============================================")
			c.Println(fmt.Sprintf("Results Saved to file: %s", filepath))
		}
		return

	},
}
