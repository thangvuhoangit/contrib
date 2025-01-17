package executor_fiber

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"logur.dev/logur"
)

func contains(needle string, slice []string) bool {
	for _, e := range slice {
		if e == needle {
			return true
		}
	}

	return false
}

func NewLoggerHandler(logger logur.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		start := time.Now()

		// Handle request, store err for logging
		chainErr := c.Next()
		if chainErr != nil {
			// Manually call error handler
			if err := c.App().ErrorHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		latency := time.Since(start)

		status := c.Response().StatusCode()

		index := 0
		switch {
		case status >= 500:
			// error index is zero
		case status >= 400:
			index = 1
		default:
			index = 2
		}

		levelIndex := index
		if levelIndex >= len(cfg.Levels) {
			levelIndex = len(cfg.Levels) - 1
		}
		level := cfg.Levels[levelIndex]

		// no log
		if level == zerolog.NoLevel || level == zerolog.Disabled {
			return nil
		}

		messageIndex := index
		if messageIndex >= len(cfg.Messages) {
			messageIndex = len(cfg.Messages) - 1
		}
		message := cfg.Messages[messageIndex]

		// logger := cfg.logger(c, latency, chainErr)

		ctx := c.UserContext()

		switch level {
		case zerolog.DebugLevel:
			logger.Debug().Ctx(ctx).Msg(message)
		case zerolog.InfoLevel:
			logger.Info().Ctx(ctx).Msg(message)
		case zerolog.WarnLevel:
			logger.Warn().Ctx(ctx).Msg(message)
		case zerolog.ErrorLevel:
			logger.Error().Ctx(ctx).Msg(message)
		case zerolog.FatalLevel:
			logger.Fatal().Ctx(ctx).Msg(message)
		case zerolog.PanicLevel:
			logger.Panic().Ctx(ctx).Msg(message)
		case zerolog.TraceLevel:
			logger.Trace().Ctx(ctx).Msg(message)
		}

		logger.Debug()

		return nil

	}
}
