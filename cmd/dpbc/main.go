package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/handlers"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
	"github.com/dracocity/draco-payment-bridge-core/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

var app *cli.App

type boundListener struct {
	network  string
	address  string
	listener net.Listener
	server   *http.Server
}

func init() {
	app = &cli.App{
		Name:  "dpbc",
		Usage: "Cryptocurrency payment gateway bridge core",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.toml",
				Usage:   "path to config TOML file (empty string for auto-detect)",
			},
		},
		Action: run,
	}
}

func run(ctx *cli.Context) error {
	configPath := ctx.String("config")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config (%s): %v", configPath, err)
	}

	if err := logger.Init(cfg.Log); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	loader := plugin.NewPluginLoader(cfg.PluginDir)
	plugins, err := loader.LoadPlugins()
	if err != nil {
		logger.Warn("failed to load plugins", "error", err)
	}
	if len(plugins) == 0 {
		logger.Warn("no payment gateways initialized. plugin directory", "plugin-dir", cfg.PluginDir)
	}

	pgRegistry := pg.NewRegistry()
	for name, p := range plugins {
		if err := p.Load(cfg.Providers[name]); err != nil {
			logger.Warn("failed to load plugin", "plugin", name, "error", err)
			continue
		}
		if ppg := p.PaymentGateway(); ppg != nil {
			pgRegistry.Register(ppg)
		} else {
			logger.Warn("plugin did not return payment gateway", "plugin", name)
		}
	}

	router := gin.New()
	router.Use(gin.Recovery())
	paymentService := services.NewPaymentService(pgRegistry)
	handler := handlers.New(paymentService)
	handler.RegisterRoutes(router)

	listeners := make([]boundListener, 0, len(cfg.Listen))
	for _, endpoint := range cfg.Listen {
		network := endpoint.Network
		address := endpoint.Address

		if network == "unix" {
			if err := os.MkdirAll(filepath.Dir(address), 0o755); err != nil {
				return err
			}
			if err := os.Remove(address); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}

		listener, err := net.Listen(network, address)
		if err != nil {
			return err
		}

		if network == "unix" {
			_ = os.Chmod(address, 0o666)
		}

		listeners = append(listeners, boundListener{
			network:  network,
			address:  address,
			listener: listener,
			server: &http.Server{
				Handler:           router,
				ReadHeaderTimeout: 5 * time.Second,
			},
		})
	}

	defer func() {
		for _, ls := range listeners {
			_ = ls.listener.Close()
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	for _, ls := range listeners {
		listener := ls
		go func() {
			logger.Info("payment bridge server listening", "network", listener.network, "address", listener.address)
			if err := listener.server.Serve(listener.listener); err != nil && err != http.ErrServerClosed {
				logger.Fatal("server error", "network", listener.network, "address", listener.address, "error", err)
			}
		}()
	}

	<-shutdownCh
	logger.Info("shutting down server")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, ls := range listeners {
		_ = ls.server.Shutdown(ctxShutdown)
		if ls.network == "unix" {
			if err := os.Remove(ls.address); err != nil && !errors.Is(err, os.ErrNotExist) {
				logger.Warn("failed to remove unix domain socket", "path", ls.address, "error", err)
			}
		}
	}

	for name, p := range plugins {
		if err := p.Unload(); err != nil {
			logger.Warn("failed to unload plugin", "plugin", name, "error", err)
		}
	}

	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
