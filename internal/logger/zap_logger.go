package logger

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	sugar  *zap.SugaredLogger
	module string
}

// WithModule returns a new logger with the specified module name.
func (l *zapLogger) withModule(module string) Logger {
	return &zapLogger{
		sugar:  l.sugar,
		module: module,
	}
}

// addModule adds the module field to the key-value pairs if it exists.
func (l *zapLogger) addModule(kv []interface{}) []interface{} {
	if l.module != "" {
		kv = append([]interface{}{"module", l.module}, kv...)
	}
	return kv
}

// Debug logs a message at the Debug level.
func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, l.addModule(keysAndValues)...)
}

// Info logs a message at the Info level.
func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, l.addModule(keysAndValues)...)
}

// Warn logs a message at the Warn level.
func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, l.addModule(keysAndValues)...)
}

// Error logs a message at the Error level.
func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, l.addModule(keysAndValues)...)
}

// Fatal logs a message at the Fatal level and then calls os.Exit(1).
func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, l.addModule(keysAndValues)...)
}

// Debugf logs a formatted string at the Debug level.
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Infof logs a formatted string at the Info level.
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warnf logs a formatted string at the Warn level.
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Errorf logs a formatted string at the Error level.
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// Fatalf logs a formatted string at the Fatal level and then calls os.Exit(1).
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}
