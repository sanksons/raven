package common

import (
	"github.com/abiosoft/ishell"
)

func GetStatusCmd() *ishell.Cmd {
	return statusCmd
}

var statusCmd = &ishell.Cmd{
	Name:     "status",
	Aliases:  []string{"st", "stats"},
	Help:     "Fetch receiver status",
	LongHelp: "",
	Func: func(c *ishell.Context) {
		rs, _ := GetRavenServer(c)
		result, err := rs.FetchStatus()
		if err != nil {
			c.Println(err.Error())
			return
		}
		type Box struct {
			Name     string
			Inflight string
			Dead     string
		}
		boxes := make([]Box, len(result.Boxes))
		for _, boxname := range result.Boxes {
			b := Box{
				Name:     boxname,
				Inflight: result.Inflight[boxname],
				Dead:     result.DeadBox[boxname],
			}
			boxes = append(boxes, b)
		}

		c.Println()
		c.Printf("Name: %s\t\t\tIsReliable: %t\n\n", result.Queue, result.IsReliable)

		c.Println("Message Box Details:")
		c.Println("----------------------------------------")
		c.Println("\tInflight\tDead\t")
		for _, box := range boxes {

			c.Printf("%s\t%s\t%s\t\n", box.Name, box.Inflight, box.Dead)
		}
		c.Println("----------------------------------------")
		c.Println()

		return
	},
}
