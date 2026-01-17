package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/bridge"
	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/handlers"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

var app *cli.App

func init() {
	app = &cli.App{
		Name:  "dpbc",
		Usage: "dpbc",
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
	// Load config
	configPath := ctx.String("config")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config (%s): %v", configPath, err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Log); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Load plugins from configured directory
	loader := plugin.NewPluginLoader(cfg.PluginDir)
	plugins, err := loader.LoadPlugins()
	if err != nil {
		logger.Warn("failed to load plugins", "error", err)
	}
	if len(plugins) == 0 {
		logger.Warn("no payment bridges initialized. plugin directory", "plugin-dir", cfg.PluginDir)
	}

	registry := bridge.NewRegistry()
	for name, p := range plugins {
		if err := p.New(); err != nil {
			logger.Warn("failed to initialize plugin", "plugin", name, "error", err)
			continue
		}
		if err := p.Load(); err != nil {
			logger.Warn("failed to load plugin", "plugin", name, "error", err)
			continue
		}
		if b := p.Bridge(); b != nil {
			registry.Register(b)
		} else {
			logger.Warn("plugin did not return bridge", "plugin", name)
		}
	}

	router := gin.New()
	router.Use(gin.Recovery())
	handler := handlers.New(registry)
	handler.RegisterRoutes(router)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("payment bridge server listening", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", "error", err)
		}
	}()

	<-shutdownCh
	logger.Info("shutting down server")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctxShutdown)

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
