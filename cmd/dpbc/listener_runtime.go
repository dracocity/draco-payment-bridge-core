package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
)

type boundListener struct {
	network  string
	address  string
	listener net.Listener
	server   *http.Server
}

func bindListeners(endpoints []config.ListenConfig, handler http.Handler) ([]boundListener, error) {
	listeners := make([]boundListener, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint.Network == "unix" {
			if err := prepareUnixSocket(endpoint.Address); err != nil {
				return nil, err
			}
		}

		listener, err := net.Listen(endpoint.Network, endpoint.Address)
		if err != nil {
			return nil, err
		}

		if endpoint.Network == "unix" {
			_ = os.Chmod(endpoint.Address, 0o666)
		}

		listeners = append(listeners, boundListener{
			network:  endpoint.Network,
			address:  endpoint.Address,
			listener: listener,
			server: &http.Server{
				Handler:           handler,
				ReadHeaderTimeout: 5 * time.Second,
			},
		})
	}

	return listeners, nil
}

func startListeners(listeners []boundListener) {
	for _, ls := range listeners {
		listener := ls
		go func() {
			logger.Info("payment bridge server listening", "network", listener.network, "address", listener.address)
			if err := listener.server.Serve(listener.listener); err != nil && err != http.ErrServerClosed {
				logger.Fatal("server error", "network", listener.network, "address", listener.address, "error", err)
			}
		}()
	}
}

func shutdownListeners(ctx context.Context, listeners []boundListener) {
	for _, ls := range listeners {
		_ = ls.server.Shutdown(ctx)
		if ls.network == "unix" {
			if err := cleanupUnixSocket(ls.address); err != nil {
				logger.Warn("failed to remove unix domain socket", "path", ls.address, "error", err)
			}
		}
	}
}

func closeListeners(listeners []boundListener) {
	for _, ls := range listeners {
		_ = ls.listener.Close()
	}
}

func prepareUnixSocket(address string) error {
	if err := os.MkdirAll(filepath.Dir(address), 0o755); err != nil {
		return err
	}

	if err := os.Remove(address); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func cleanupUnixSocket(address string) error {
	if err := os.Remove(address); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}
