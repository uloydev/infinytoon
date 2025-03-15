package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("GET / endpoint hit")
		return c.SendString("Hello, World!")
	})

	// test endpoint
	app.Get("/test", func(c *fiber.Ctx) error {
		log.Println("GET /test endpoint hit")
		return c.SendString("Hello, Test!")
	})

	// post endpoint
	app.Post("/post", func(c *fiber.Ctx) error {
		log.Println("POST /post endpoint hit")
		// return the body of the post request
		return c.Send(c.Body())
	})

	app.Listen(":3000")

}
