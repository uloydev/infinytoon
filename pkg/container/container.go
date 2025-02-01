package container

import (
	"context"
	"log"
	"sync"

	"infinitoon.dev/infinitoon/pkg/cmd"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

type Container struct {
	appCtx   *appctx.AppContext
	commands map[string]cmd.Command
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewContainer(appCtx *appctx.AppContext) *Container {
	ctx, cancel := context.WithCancel(context.Background())
	return &Container{
		appCtx:   appCtx,
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
		log.Println("Starting command:", name)
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
		log.Println("Shutting down command:", name)
		if err := command.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}
