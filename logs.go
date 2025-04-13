package fibber

import (
	"fmt"
	"strings"

	config "github.com/robbyriverside/fibber/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigPath() string {
	return config.Path()
}

func InitLogger(env string) {
	once.Do(func() {
		env = strings.ToLower(env)
		if env == "" || env == "production" {
			env = "production"
		} else if env == "dev" {
			env = "development"
		}
		Options.Environment = env

		cfg := zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.LevelKey = "level"
		cfg.EncoderConfig.MessageKey = "msg"
		cfg.EncoderConfig.CallerKey = "caller"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

		log, err := cfg.Build(zap.Fields(
			zap.String("app", Options.AppName),
			zap.String("version", Options.Version),
			zap.String("env", Options.Environment),
		))
		if err != nil {
			panic(err)
		}
		logger = log.Sugar()
	})
}

// VLogf logs only when verbose is true
func VLogf(format string, args ...interface{}) {
	if Options.Verbose {
		fmt.Printf("[verbose] "+format+"\n", args...)
	}
}
