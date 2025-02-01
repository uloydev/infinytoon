package rest

import "github.com/gofiber/fiber/v2"

type RestRoute struct {
	prefix string
	isRoot bool
	handle func(router fiber.Router)
}

func NewRestRoute() *RestRoute {
	return &RestRoute{}
}

func (r *RestRoute) SetRoot() *RestRoute {
	r.isRoot = true
	return r
}

func (r *RestRoute) SetPrefix(prefix string) *RestRoute {
	r.prefix = prefix
	return r
}

func (r *RestRoute) Handler(handle func(router fiber.Router)) *RestRoute {
	r.handle = handle
	return r
}

func (r *RestRoute) Register(basePath string, router fiber.Router) {
	if !r.isRoot {
		router = router.Group(basePath)
		if r.prefix != "" {
			router = router.Group(r.prefix)
		}
	}

	r.handle(router)
}
