package executor_fiber

import (
	"context"

	"github.com/gofiber/fiber/v2"
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
	PreHooks    []func(ctx context.Context)
	PostHooks   []func(err error)
}

type FiberExecutorOption func(*fiber.App)

func configDefault(config FiberExecutorConfig) FiberExecutorConfig {
	var cfg FiberExecutorConfig
	cfg.FiberConfig = defaultFiberConfig(config.FiberConfig)
	cfg.Options = config.Options
	cfg.FiberConfig.AppName = config.Name

	return cfg
}

func NewFiberExecutor(config FiberExecutorConfig) *gcakit.Executor {
	cfg := configDefault(config)
	engine := fiber.New(cfg.FiberConfig)

	for _, hook := range cfg.Options {
		hook(engine)
	}

	executor := &gcakit.Executor{
		Name: cfg.Name,
		Execute: func(ctx context.Context) error {
			for _, prehook := range cfg.PreHooks {
				prehook(ctx)
			}

			return engine.Server().ListenAndServe(cfg.Addr)
		},
		Interrupt: func(err error) {
			for _, posthook := range cfg.PostHooks {
				posthook(err)
			}

			engine.Shutdown()
		},
	}

	// engine.Use(fiberzap.New(fiberzap.Config{
	// 	Logger:   zapLogger,
	// 	Fields:   []string{"protocol", "pid", "ip", "host", "url", "route", "method", "queryParams", "bytesReceived", "bytesSent"},
	// 	SkipURIs: []string{"/swagger/*"},
	// }))

	// engine.Use(requestid.New())
	// engine.Use(recover.New())
	// engine.Use(cors.New())
	// engine.Use(compress.New(compress.Config{
	// 	Level: compress.LevelBestSpeed,
	// }))
	// engine.Use(helmet.New())

	// redisStorage := redis.New(redis.Config{
	// 	Host:      cfg.Redis.Host,
	// 	Port:      cfg.Redis.Port,
	// 	Username:  cfg.Redis.Username,
	// 	Password:  cfg.Redis.Password,
	// 	Database:  cfg.Redis.DB,
	// 	Reset:     false,
	// 	TLSConfig: nil,
	// 	PoolSize:  10 * runtime.GOMAXPROCS(0),
	// })

	// fiberSession := sessions.NewFiberSession(redisStorage)

	// endpoint.ConfigureEndpoints(ctx, engine, logger, cfg, fiberSession)
	// endpoint.ConfigureSwagger(ctx, engine, logger, cfg)

	// ctx.AddToStartUpMessage("HTTP server", fmt.Sprintf("http://127.0.0.1%s (bound on host 0.0.0.0 and port %s)", cfg.HTTP.Port, cfg.HTTP.Port))

	// return func() { engine.Shutdown() }, engine.Listen(cfg.HTTP.Port)

	return executor
}
