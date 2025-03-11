package container

import (
	"context"
	"sync"

	"infinitoon.dev/infinitoon/pkg/cmd"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
)

type Container struct {
	appCtx   *appctx.AppContext
	log      *logger.Logger
	commands map[string]cmd.Command
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewContainer(appCtx *appctx.AppContext) *Container {
	ctx, cancel := context.WithCancel(context.Background())
	return &Container{
		appCtx:   appCtx,
		log:      appCtx.Get(appctx.LoggerKey).(*logger.Logger),
		commands: make(map[string]cmd.Command),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (c *Container) RegisterCommand(commands ...cmd.Command) {
	for _, command := range commands {
		c.commands[command.Name()] = command
	}
}

func (c *Container) Run() error {
	for name, command := range c.commands {
		c.log.Info().Msgf("Starting command: %s", name)
		c.wg.Add(1)
		go func(cmd cmd.Command) {
			defer c.wg.Done()
			if err := cmd.Run(); err != nil {
				c.cancel()
			}
		}(command)
	}
	c.wg.Wait()
	return c.ctx.Err()
}

func (c *Container) Shutdown() error {
	c.cancel()
	for name, command := range c.commands {
		c.log.Info().Msgf("Shutting down command: %s", name)
		if err := command.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}
