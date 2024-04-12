package common

import (
	"github.com/kloudlite/iot-devices/pkg/k3s"
	"github.com/kloudlite/iot-devices/pkg/logging"
)

/*
steps to implement:
[ ] setup k3s server
*/
func reconDevice(l logging.Logger) error {
	c := k3s.New(l)

	if err := c.CheckIfInstalled(); err != nil {
		if err := c.Install(); err != nil {
			return err
		}
	}

	defer func() {
		if err := c.Uninstall(); err != nil {
			return
		}
	}()

	return nil
}
