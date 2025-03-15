package controller

import (
	"infinitoon.dev/infinitoon/apps/proxy/config"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/rest"
)

type IController interface {
	Route() *rest.RestRoute
}

type Controller struct {
	appCtx *appctx.AppContext

	ctrls []IController
}

func InitController(appCtx *appctx.AppContext, cfg *config.Config) *Controller {
	return &Controller{
		appCtx: appCtx,
		ctrls: []IController{
			// register all controllers here
			NewRootController(appCtx, cfg),
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
