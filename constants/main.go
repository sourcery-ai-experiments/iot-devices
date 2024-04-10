package constants

import "fmt"

const (
	IotServerEndpoint = "iotnet.dev.kloudlite.io"

	BroadcastIP   = "255.255.255.255"
	BroadcastPort = 12345

	BroadcastInterval = 4

	ExpireDuration = 10

	ProxyServerPort = 8000
)

func GetHealthyUrl() string {
	return fmt.Sprintf("https://%s/healthy", IotServerEndpoint)
}

func GetIotPingPath() string {
	return "/healthy"
}

func GetIotServerEndpoint() string {
	return fmt.Sprintf("https://%s", IotServerEndpoint)
}
