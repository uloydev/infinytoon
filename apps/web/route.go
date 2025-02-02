package main

import (
	"infinitoon.dev/infinitoon/apps/web/controller"
	"infinitoon.dev/infinitoon/apps/web/repository"
	"infinitoon.dev/infinitoon/apps/web/service"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/rest"
)

func GetRoutes(appCtx *appctx.AppContext) []*rest.RestRoute {
	// repositories
	userRepo := repository.NewUserRepo(appCtx)

	// services
	userService := service.NewUserService(appCtx, userRepo)

	// controllers
	userController := controller.NewUserController(appCtx, userService)

	return []*rest.RestRoute{
		userController.Route(),
	}
}
