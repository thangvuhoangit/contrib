package executor_fiber

import (
	"context"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/thangvuhoangit/gcakit"
)

var FiberConfigDefault = fiber.Config{
	BodyLimit:             4 * 1024 * 1024,
	Concurrency:           256 * 1024,
	AppName:               "fiber http server",
	DisableStartupMessage: true,
}

func defaultFiberConfig(fiberCfg ...fiber.Config) fiber.Config {
	if len(fiberCfg) < 1 {
		return FiberConfigDefault
	}

	// Override default config
	cfg := fiberCfg[0]
	return cfg
}

type FiberExecutorConfig struct {
	Addr        string // ":8080" or "127.0.0.1:8080"
	Name        string // "fiber http server"
	FiberConfig fiber.Config
	Options     []FiberExecutorOption
	PreHooks    []func(ctx context.Context, engine *fiber.App)
	PostHooks   []func(err error, engine *fiber.App)
}

type FiberExecutorOption func(*fiber.App)

func configDefault(config FiberExecutorConfig) FiberExecutorConfig {
	var cfg FiberExecutorConfig
	cfg.FiberConfig = defaultFiberConfig(config.FiberConfig)
	cfg.Options = config.Options
	cfg.FiberConfig.AppName = config.Name
	cfg.PreHooks = config.PreHooks
	cfg.PostHooks = config.PostHooks
	cfg.Name = config.Name
	cfg.Addr = config.Addr

	return cfg
}

func NewFiberExecutor(config FiberExecutorConfig) *gcakit.Executor {
	cfg := configDefault(config)
	engine := fiber.New(cfg.FiberConfig)
	log.SetOutput(os.Stdout)

	if len(cfg.Options) > 0 {
		owg := new(sync.WaitGroup)
		owg.Add(len(cfg.Options))
		for _, hook := range cfg.Options {
			go func() {
				hook(engine)
				owg.Done()
			}()
		}
		owg.Wait()
	}

	executor := &gcakit.Executor{
		Name: cfg.Name,
		Execute: func(ctx context.Context) error {
			if len(cfg.PreHooks) > 0 {
				wg := new(sync.WaitGroup)
				wg.Add(len(cfg.PreHooks))
				for _, prehook := range cfg.PreHooks {
					go func() {
						prehook(ctx, engine)
						wg.Done()
					}()
				}
				wg.Wait()
			}

			go engine.Listen(cfg.Addr)
			return nil
		},
		Interrupt: func(err error) {
			wg := new(sync.WaitGroup)
			wg.Add(len(cfg.PostHooks))
			for _, posthook := range cfg.PostHooks {
				go func() {
					posthook(err, engine)
					wg.Done()
				}()
			}
			wg.Wait()
			engine.ShutdownWithContext(context.Background())
		},
	}

	return executor
}
