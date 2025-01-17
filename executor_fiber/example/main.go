package main

import (
	"context"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/thangvuhoangit/contrib/executor_fiber"
	"github.com/thangvuhoangit/gcakit"
	"go.uber.org/zap"
)

func main() {
	myApp := gcakit.New(gcakit.Config{Name: "My App"})
	logger := myApp.Logger()
	zapLogger, _ := zap.NewProduction()

	myApp.AddExecutor(executor_fiber.NewFiberExecutor(executor_fiber.FiberExecutorConfig{
		Name: "MyFiberExecutor",
		Addr: ":5000",
		FiberConfig: fiber.Config{
			DisableStartupMessage: true,
		},
		PreHooks: []func(ctx context.Context, engine *fiber.App){
			func(ctx context.Context, engine *fiber.App) {
				// fmt.Println("PreHook 1")
				engine.Use(fiberzap.New(fiberzap.Config{
					Logger:   zapLogger,
					Fields:   []string{"protocol", "pid", "ip", "host", "url", "route", "method", "queryParams", "bytesReceived", "bytesSent"},
					SkipURIs: []string{"/swagger/*"},
				}))

				engine.Use(executor_fiber.NewLoggerHandler(logger))
			},
			func(ctx context.Context, engine *fiber.App) {
				// fmt.Println("PreHook 2")
				engine.Use(requestid.New())
				engine.Use(recover.New())
				engine.Use(cors.New())
				engine.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
				engine.Use(helmet.New())
			},
			func(ctx context.Context, engine *fiber.App) {
				// fmt.Println("PreHook 3")
				v1 := engine.Group("/api/v1")
				v1.Get("/hello", func(c *fiber.Ctx) error {
					return c.SendString("Hello, World 👋!")
				})
			},
		},
	}))

	myApp.Start()

	<-myApp.WaitDone()

	myApp.Stop()
}
