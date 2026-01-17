package logger

import (
	"os"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var gZapLogger Logger

// Init initializes the global zap logger based on the provided configuration.
func Init(cfg config.LogConfig) error {
	// Set the logging level from the configuration.
	// If the level is invalid, it defaults to InfoLevel.
	var level zapcore.Level
	if err := level.Set(cfg.Level); err != nil {
		level = zapcore.InfoLevel
	}

	// Configure the log encoder (e.g., JSON or console format).
	encoderCfg := zap.NewProductionEncoderConfig()
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	// Configure the log output destinations (stdout, stderr, file).
	var writeSyncers []zapcore.WriteSyncer
	if cfg.Output.Stdout {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	if cfg.Output.Stderr {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stderr))
	}
	if cfg.Output.FilePath != "" {
		writeSyncers = append(writeSyncers, zapcore.AddSync(&lumberjack.Logger{
			Filename: cfg.Output.FilePath,
			MaxSize:  cfg.Rotation.MaxSize,
			MaxAge:   cfg.Rotation.MaxAge,
			Compress: cfg.Rotation.Compress,
		}))
	}
	// Combine multiple outputs into a single WriteSyncer.
	multiWS := zapcore.NewMultiWriteSyncer(writeSyncers...)

	// Create the core logger with the encoder, writer, and level.
	core := zapcore.NewCore(encoder, multiWS, level)
	// Create the final zap logger, adding the caller information to log entries.
	zapLog := zap.New(core, zap.AddCaller())

	// Set the global logger instance with a sugared logger for easier use.
	gZapLogger = &zapLogger{
		sugar:  zapLog.Sugar(),
		module: "global",
	}
	return nil
}

// Sync flushes any buffered log entries.
func Sync() error {
	return gZapLogger.(*zapLogger).sugar.Sync()
}

// getZapLogger returns the global zap logger instance.
// It panics if the logger has not been initialized.
func getZapLogger() *zapLogger {
	if gZapLogger == nil {
		panic("logger: global zap logger is not initialized. Call Init() first")
	}
	return gZapLogger.(*zapLogger)
}

// WithModule returns a new logger with the specified module name.
func WithModule(module string) Logger {
	return getZapLogger().withModule(module)
}

// Debug logs a message at the Debug level.
func Debug(msg string, keysAndValues ...interface{}) {
	l := getZapLogger()
	l.sugar.Debugw(msg, l.addModule(keysAndValues)...)
}

// Info logs a message at the Info level.
func Info(msg string, keysAndValues ...interface{}) {
	l := getZapLogger()
	l.sugar.Infow(msg, l.addModule(keysAndValues)...)
}

// Warn logs a message at the Warn level.
func Warn(msg string, keysAndValues ...interface{}) {
	l := getZapLogger()
	l.sugar.Warnw(msg, l.addModule(keysAndValues)...)
}

// Error logs a message at the Error level.
func Error(msg string, keysAndValues ...interface{}) {
	l := getZapLogger()
	l.sugar.Errorw(msg, l.addModule(keysAndValues)...)
}

// Fatal logs a message at the Fatal level and then calls os.Exit(1).
func Fatal(msg string, keysAndValues ...interface{}) {
	l := getZapLogger()
	l.sugar.Fatalw(msg, l.addModule(keysAndValues)...)
}

// Debugf logs a formatted string at the Debug level.
func Debugf(format string, args ...interface{}) {
	l := getZapLogger()
	l.sugar.Debugf(format, args...)
}

// Infof logs a formatted string at the Info level.
func Infof(format string, args ...interface{}) {
	l := getZapLogger()
	l.sugar.Infof(format, args...)
}

// Warnf logs a formatted string at the Warn level.
func Warnf(format string, args ...interface{}) {
	l := getZapLogger()
	l.sugar.Warnf(format, args...)
}

// Errorf logs a formatted string at the Error level.
func Errorf(format string, args ...interface{}) {
	l := getZapLogger()
	l.sugar.Errorf(format, args...)
}

// Fatalf logs a formatted string at the Fatal level and then calls os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	l := getZapLogger()
	l.sugar.Fatalf(format, args...)
}
