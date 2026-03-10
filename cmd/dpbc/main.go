package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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
	cfg, err := loadConfig(ctx.String("config"))
	if err != nil {
		return err
	}

	if err := initLogger(cfg); err != nil {
		return err
	}
	defer logger.Sync()

	plugins, registry := loadPaymentGateways(cfg)
	defer unloadPlugins(plugins)

	router := newRouter(registry)

	runtime, err := newServerRuntime(cfg.Listen, router)
	if err != nil {
		return err
	}
	defer runtime.close()

	runtime.start()
	waitForShutdownSignal()
	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return runtime.shutdown(shutdownCtx)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}

func loadConfig(path string) (*config.Config, error) {
	cfg, err := config.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config (%s): %w", path, err)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *config.Config) error {
	if len(cfg.Listen) == 0 {
		return errors.New("invalid config: listen endpoints are required")
	}
	for i, endpoint := range cfg.Listen {
		if strings.TrimSpace(endpoint.Network) == "" {
			return fmt.Errorf("invalid config: listen[%d].network is required", i)
		}
		if strings.TrimSpace(endpoint.Address) == "" {
			return fmt.Errorf("invalid config: listen[%d].address is required", i)
		}
	}

	return nil
}

func initLogger(cfg *config.Config) error {
	if err := logger.Init(cfg.Log); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	return nil
}

func loadPaymentGateways(cfg *config.Config) (map[string]plugin.Plugin, *pg.Registry) {
	loader := plugin.NewPluginLoader(cfg.PluginDir)
	plugins, err := loader.LoadPlugins()
	if err != nil {
		logger.Warn("failed to load plugins", "error", err)
	}
	if len(plugins) == 0 {
		logger.Warn("no payment gateways initialized. plugin directory", "plugin-dir", cfg.PluginDir)
	}

	registry := pg.NewRegistry()
	for name, p := range plugins {
		if err := p.Load(cfg.Providers[name]); err != nil {
			logger.Warn("failed to load plugin", "plugin", name, "error", err)
			continue
		}
		if gateway := p.PaymentGateway(); gateway != nil {
			registry.Register(gateway)
			continue
		}
		logger.Warn("plugin did not return payment gateway", "plugin", name)
	}

	return plugins, registry
}

func unloadPlugins(plugins map[string]plugin.Plugin) {
	for name, p := range plugins {
		if err := p.Unload(); err != nil {
			logger.Warn("failed to unload plugin", "plugin", name, "error", err)
		}
	}
}

func newRouter(registry *pg.Registry) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	paymentService := services.NewPaymentService(registry)
	handlers.New(paymentService).RegisterRoutes(router)
	return router
}

func waitForShutdownSignal() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(shutdownCh)
	<-shutdownCh
}

type serverRuntime struct {
	listeners []boundListener
}

func newServerRuntime(endpoints []config.ListenConfig, handler http.Handler) (*serverRuntime, error) {
	listeners, err := bindListeners(endpoints, handler)
	if err != nil {
		return nil, err
	}
	return &serverRuntime{listeners: listeners}, nil
}

func (s *serverRuntime) start() {
	for _, ls := range s.listeners {
		listener := ls
		go func() {
			logger.Info("payment bridge server listening", "network", listener.network, "address", listener.address)
			if err := listener.server.Serve(listener.listener); err != nil && err != http.ErrServerClosed {
				logger.Fatal("server error", "network", listener.network, "address", listener.address, "error", err)
			}
		}()
	}
}

func (s *serverRuntime) shutdown(ctx context.Context) error {
	var firstErr error
	for _, ls := range s.listeners {
		if err := ls.server.Shutdown(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
		if ls.network == "unix" {
			if err := cleanupUnixSocket(ls.address); err != nil {
				logger.Warn("failed to remove unix domain socket", "path", ls.address, "error", err)
			}
		}
	}
	return firstErr
}

func (s *serverRuntime) close() {
	for _, ls := range s.listeners {
		_ = ls.listener.Close()
	}
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
