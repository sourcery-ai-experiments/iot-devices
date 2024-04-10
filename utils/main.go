package utils

import (
	"net/http"
	"time"

	"github.com/kloudlite/iot-devices/constants"
)

func IsConn() bool {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(constants.GetHealthyUrl())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}
