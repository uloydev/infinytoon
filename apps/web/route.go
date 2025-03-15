package main

import (
	"github.com/gofiber/fiber/v2"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/rest"
)

func GetRoutes(appCtx *appctx.AppContext) []*rest.RestRoute {
	// repositories
	// userRepo := repository.NewUserRepo(appCtx)

	// // services
	// userService := service.NewUserService(appCtx, userRepo)

	// // controllers
	// userController := controller.NewUserController(appCtx, userService)

	return []*rest.RestRoute{
		rest.NewRestRoute().Handler(func(router fiber.Router) {
			router.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("Hello, World!")
			})
		}),
	}
}
