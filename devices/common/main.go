package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	// "github.com/kloudlite/iot-devices/pkg/conf"
	"github.com/kloudlite/iot-devices/pkg/logging"
)

type Response struct {
	AccountName       string `json:"accountName"`
	CreationTime      string `json:"creationTime"`
	DeploymentName    string `json:"deploymentName"`
	DisplayName       string `json:"displayName"`
	ID                string `json:"id"`
	IP                string `json:"ip"`
	MarkedForDeletion string `json:"markedForDeletion"`
	Name              string `json:"name"`
	PodCIDR           string `json:"podCIDR"`
	ProjectName       string `json:"projectName"`
	PublicKey         string `json:"publicKey"`
	RecordVersion     int    `json:"recordVersion"`
	ServiceCIDR       string `json:"serviceCIDR"`
	UpdateTime        string `json:"updateTime"`
	Version           string `json:"version"`
}

func (c *Response) FromJson(data []byte) error {
	return json.Unmarshal(data, c)
}

var (
	domains = []string{
		constants.IotServerEndpoint,
		constants.DnsDomain,
		"get.k3s.io",
		"ghcr.io",
		"registry.hub.docker.com",
	}

	device *Response = nil
)

func GetDomains() []string {
	return domains
}

func GetDevice() (*Response, error) {
	if device == nil {
		return nil, fmt.Errorf("device is not initialized")
	}
	return device, nil
}

func StartPing(l logging.Logger) {
	for {
		if err := ping(l); err != nil {
			l.Errorf(err, "sending ping to server")
		}

		time.Sleep(constants.PingInterval * time.Second)
	}
}

func ping(l logging.Logger) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// c, err := conf.GetConf()
	// if err != nil {
	// 	return err
	// }

	var data = struct {
		PublicKey string `json:"publicKey"`
	}{
		// PublicKey: c.PublicKey,
		PublicKey: "10.2.2.2",
	}

	dataBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	resp, err := client.Post(constants.GetPingUrl(), "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response Response

	// read all the response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	l.Infof("Ping response: %s", buf.String())
	if resp.StatusCode == http.StatusOK {

		if err := response.FromJson(buf.Bytes()); err != nil {
			return err
		}

		device = &response

		if err := reconDevice(l); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("status code: %d", resp.StatusCode)
}
