package k3s

import (
	"github.com/kloudlite/iot-devices/pkg/logging"
	"github.com/kloudlite/iot-devices/utils"
)

type client struct {
	l logging.Logger
}

func New(l logging.Logger) *client {
	return &client{
		l: l,
	}
}

func (c *client) CheckIfInstalled() error {
	if err := utils.ExecCmd("k3s --version", true); err != nil {
		return err
	}

	if err := utils.ExecCmd("systemctl status k3s", true); err != nil {
		return err
	}

	return nil
}

func (c *client) Install() error {

	if err := utils.ExecCmd("curl -sfL https://get.k3s.io | sh -", true); err != nil {
		return err
	}

	return nil
}

func (c *client) Uninstall() error {
	if err := utils.ExecCmd("k3s-uninstall.sh", true); err != nil {
		return err
	}

	return nil
}
