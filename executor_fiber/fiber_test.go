package executor_fiber

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFiberConfig(t *testing.T) {
	config := defaultFiberConfig()
	assert.Equal(t, FiberConfigDefault, config)

	customConfig := fiber.Config{AppName: "custom app"}
	config = defaultFiberConfig(customConfig)
	assert.Equal(t, "custom app", config.AppName)
}

func TestConfigDefault(t *testing.T) {
	config := FiberExecutorConfig{
		Name: "test server",
	}
	cfg := configDefault(config)
	assert.Equal(t, "test server", cfg.FiberConfig.AppName)
}

// func TestNewFiberExecutor(t *testing.T) {
// 	config := FiberExecutorConfig{
// 		Addr: ":8080",
// 		Name: "test server",
// 	}
// 	executor := NewFiberExecutor(config)
// 	assert.NotNil(t, executor)
// 	assert.Equal(t, "test server", executor.Name)

// 	ctx := context.Background()
// 	err := executor.Execute(ctx)
// 	assert.Nil(t, err)

// 	executor.Interrupt(nil)
// }
