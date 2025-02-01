package rest_middleware

import (
	"github.com/gofiber/fiber/v2"
	"infinitoon.dev/infinitoon/pkg/utils"
)

// func to return fiber handler to set X-REQUEST-ID to response header
func RequestIDMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		c.Set("X-INFINITOON-REQUEST-ID", utils.GetUUID())
		return c.Next()
	}

}
