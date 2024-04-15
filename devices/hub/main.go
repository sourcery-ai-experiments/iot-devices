package hub

import (
	"github.com/kloudlite/api/pkg/logging"
	"github.com/kloudlite/iot-devices/types"
)

type client struct {
	ctx    types.MainCtx
	logger logging.Logger

	domains Dms
}

func Run(ctx types.MainCtx) error {
	c := &client{
		ctx:     ctx,
		logger:  ctx.GetLogger(),
		domains: Dms{},
	}

	c.logger.Infof("Starting hub")

	go c.resyncDomains()

	if err := c.selfBroadcast(); err != nil {
		c.logger.Errorf(err, "Error broadcasting message")

		// TODO: handle error
		return err
	}
	return nil
}
