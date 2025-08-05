package infra

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var logger *zap.Logger
var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    encodeLevel(),
	EncodeTime:     zapcore.RFC3339TimeEncoder,
	EncodeDuration: zapcore.MillisDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func GetLogger() *zap.Logger {
	if logger == nil {
		loggerCfg := &zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		var err error
		logger, err = loggerCfg.Build(zap.AddStacktrace(zap.DPanicLevel))
		if err != nil {
			log.Println("Error creating logger:", err)
			logger, err = zap.NewProduction()
			if err != nil {
				log.Println("Error creating default logger", err)
				logger = zap.NewNop()
			}
		}
	}
	return logger
}

// SetLogger Make the singleton opt-in for tests
func SetLogger(l *zap.Logger) {
	if logger != nil {
		err := logger.Sync()
		if err != nil {
			log.Println("Error syncing logger:", err)
		}
	}
	logger = l
}

func DestroyLogger() {
	if logger != nil {
		err := logger.Sync()
		if err != nil {
			log.Println("Error syncing logger:", err)
		}
		logger = nil
	}
}

func encodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		}
	}
}
