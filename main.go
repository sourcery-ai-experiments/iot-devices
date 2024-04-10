package main

import (
	"context"
	"time"

	"github.com/kloudlite/iot-devices/devices/hub"
	"github.com/kloudlite/iot-devices/devices/local"
	"github.com/kloudlite/iot-devices/pkg/logging"
	"github.com/kloudlite/iot-devices/utils"
)

func main() {
	// if err := Run(); err != nil {
	// 	println(err.Error())
	// }

	if err := OnlyLocal(); err != nil {
		println(err.Error())
	}
	//
	// if err := OnlyHub(); err != nil {
	// 	println(err.Error())
	// }
}

func OnlyLocal() error {
	l, err := logging.New(&logging.Options{})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	println("Connection is unhealthy")
	if err := local.Run(ctx, l); err != nil {
		l.Errorf(err, "Error running local")
		return err
	}

	return nil
}

func OnlyHub() error {
	l, err := logging.New(&logging.Options{})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := hub.Run(ctx, l); err != nil {
		l.Errorf(err, "Error running hub")
		return err
	}

	return nil
}

func Run() error {
	l, err := logging.New(&logging.Options{})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	var isConnected = utils.IsConn()

	go func() {
		for {
			ic := utils.IsConn()
			if isConnected != ic {
				isConnected = ic
				cancel()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	for {
		ctx, cancel = context.WithCancel(context.Background())
		if isConnected {
			if err := hub.Run(ctx, l); err != nil {
				l.Errorf(err, "Error running hub")
				continue
			}
		}

		println("Connection is unhealthy")
		if err := local.Run(ctx, l); err != nil {
			l.Errorf(err, "Error running local")
			continue
		}
	}

}
