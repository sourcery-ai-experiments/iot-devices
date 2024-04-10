package local

import (
	"context"
	"fmt"
	"time"

	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/kloudlite/iot-devices/constants"
)

var (
	currentTargetIndex int32 = 0
)

func (c *client) listenProxy() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := c.getNextTarget()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse target URL: %s", err.Error()), http.StatusServiceUnavailable)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		r.Host = targetURL.Host
		proxy.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", constants.ProxyServerPort),
		Handler: mux,
	}

	// Listen for context cancellation in a separate goroutine
	go func() {
		<-c.ctx.Done() // This blocks until the context is cancelled

		// Context has been cancelled, shutdown the server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			c.logger.Errorf(err, "Failed to gracefully shutdown server")
		}
	}()

	c.logger.Infof("Starting server on port %d", constants.ProxyServerPort)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// ListenAndServe always returns a non-nil error. ErrServerClosed is returned when we call Shutdown
		c.logger.Errorf(err, "Failed to start server")
		return err
	}

	return nil
}

func (c *client) getNextTarget() (*url.URL, error) {
	targets := hubs.GetHubs()

	if len(targets) == 0 {
		return nil, fmt.Errorf("No targets available")
	}

	index := atomic.AddInt32(&currentTargetIndex, 1)
	index = index % int32(len(targets))

	targetURL, err := url.Parse(fmt.Sprintf("http://%s:%d", targets[index], constants.ProxyServerPort))
	if err != nil {
		return nil, err
	}
	return targetURL, nil
}
