package hub

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/kloudlite/iot-devices/constants"
)

func (c *client) listenProxy() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(constants.GetIotServerEndpoint())
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
