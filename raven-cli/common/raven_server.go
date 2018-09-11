package common

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
)

type RecStatus struct {
	Queue      string
	Boxes      []string
	DeadBox    map[string]string
	Inflight   map[string]string
	IsReliable bool
}

type RavenServer struct {
	IP   string
	Port string
}

func (this *RavenServer) Ping() (string, error) {
	url := fmt.Sprintf("http://%s:%s/ping", this.IP, this.Port)
	var data string
	err := fireHttp(url, &data)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (this *RavenServer) FetchStatus() (RecStatus, error) {
	url := fmt.Sprintf("http://%s:%s/stats", this.IP, this.Port)
	var data RecStatus
	err := fireHttp(url, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *RavenServer) FlushDeadQ() (map[string]interface{}, error) {
	url := fmt.Sprintf("http://%s:%s/flushDead", this.IP, this.Port)
	var data map[string]interface{}
	err := fireHttpPost(url, "", &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *RavenServer) FlushAll() (map[string]interface{}, error) {
	url := fmt.Sprintf("http://%s:%s/flushAll", this.IP, this.Port)
	var data map[string]interface{}
	err := fireHttpPost(url, "", &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *RavenServer) GetNewShell() *ishell.Shell {
	ishell := ishell.NewWithConfig(
		&readline.Config{Prompt: fmt.Sprintf("Raven@%s:%s >>> ", this.IP, this.Port)},
	)
	ishell.Set("raven", *this)
	return ishell
}
