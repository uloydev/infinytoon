package main

import (
	"log"

	"infinitoon.dev/infinitoon/apps/web/middleware"
	"infinitoon.dev/infinitoon/pkg/cmd"
	"infinitoon.dev/infinitoon/pkg/container"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/database"
)

func main() {
	appCtx := appctx.NewAppContext()
	db := database.NewMongoDB(appCtx, "mongodb://root:password@mongodb:27017/infinitoon?authSource=admin")
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	container := container.NewContainer(appCtx)
	container.RegisterCommand(
		cmd.NewRestCommand(appCtx, &cmd.RestCommandConfig{
			Name:        "infinitoon website",
			Host:        "0.0.0.0",
			Port:        "8080",
			BasePath:    "/",
			Middlewares: middleware.GetMiddlewares(),
			Routes:      GetRoutes(appCtx),
		}),
	)

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}

	if err := container.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
