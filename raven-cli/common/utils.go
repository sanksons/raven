package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/abiosoft/ishell"
)

const HTTP_TIMEOUT = 60 //in seconds

func fireHttp(url string, payload interface{}) error {
	client := http.Client{Timeout: HTTP_TIMEOUT * time.Second}

	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error Reading response body, Error: %s", err.Error())
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d, Result: %s", response.StatusCode, string(bodyBytes))
	}
	err = json.Unmarshal(bodyBytes, payload)
	if err != nil {
		return fmt.Errorf("StatusCode: %d, Unmarshalling failed, Error: %s", response.StatusCode, err.Error())
	}
	return nil
}

func GetRavenServer(c *ishell.Context) (*RavenServer, error) {
	i := c.Get("raven")
	if s, ok := i.(RavenServer); ok {
		return &s, nil
	}
	return nil, fmt.Errorf("could not get server")
}
