package controller

import (
	"infinitoon.dev/infinitoon/apps/proxy/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/database"
	"infinitoon.dev/infinitoon/pkg/rest"
)

type IController interface {
	Route() *rest.RestRoute
}

type Controller struct {
	appCtx *appctx.AppContext
	kv     *database.KVClient
	ctrls  []IController
}

func InitController(appCtx *appctx.AppContext, cfg *config.Config) *Controller {
	kv := database.GetKVClientFromCtx(appCtx)
	return &Controller{
		appCtx: appCtx,
		kv:     kv,
		ctrls: []IController{
			// register all controllers here
			NewRootController(appCtx, cfg, kv),
		},
	}
}

func (c *Controller) Routes() []*rest.RestRoute {
	routes := []*rest.RestRoute{}

	for _, ctrl := range c.ctrls {
		routes = append(routes, ctrl.Route())
	}

	return routes
}
