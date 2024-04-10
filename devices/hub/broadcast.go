package hub

import (
	"fmt"
	"net"
	"time"

	"github.com/kloudlite/iot-devices/constants"
)

const udpConnType = "udp"

func (c *client) selfBroadcast() error {
	c.logger.Infof("Broadcasting message...")
	defer c.logger.Infof("Broadcasting message ended")

	broadcastAddr, err := net.ResolveUDPAddr(udpConnType, fmt.Sprintf("%s:%d", constants.BroadcastIP, constants.BroadcastPort))
	if err != nil {
		return err
	}

	localAddr, err := net.ResolveUDPAddr(udpConnType, ":0")
	if err != nil {
		return err
	}

	conn, err := net.DialUDP(udpConnType, localAddr, broadcastAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	message := []byte("Ping")

	for {

		if c.ctx.Err() != nil {
			return fmt.Errorf("Context cancelled")
		}

		_, err := conn.Write(message)
		if err != nil {
			// TODO: handle error
			return err
		}
		time.Sleep(constants.BroadcastInterval * time.Second)
	}
}
