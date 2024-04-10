package local

import (
	"context"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/pkg/logging"
)

type hub struct {
	lastPing time.Time
}

type hubstype map[string]hub

func (h *hubstype) cleanup() {
	for k, v := range *h {
		if time.Since(v.lastPing) > constants.ExpireDuration*time.Second {
			delete(*h, k)
		}
	}
}

func (h *hubstype) GetHubs() []string {
	h.cleanup()

	var hubs []string
	for k := range *h {
		hubs = append(hubs, k)
	}

	return hubs
}

var hubs = hubstype{}

type client struct {
	logger logging.Logger
	ctx    context.Context
}

func Run(ctx context.Context, logger logging.Logger) error {
	c := &client{
		logger: logger,
		ctx:    ctx,
	}

	c.logger.Infof("Starting local")

	go func() {
		if err := c.listenProxy(); err != nil {
			// TODO: handle error
			panic(err)
		}
	}()

	if err := c.listenBroadcast(); err != nil {
		// TODO: handle error
		return err
	}

	return nil
}
