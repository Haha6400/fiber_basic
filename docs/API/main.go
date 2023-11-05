package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	micro := fiber.New()

	/*
		Route Handlers
	*/
	// Simple GET handler
	app.Get("/api/list", func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request!")
	})

	// Simple POST handler
	app.Post("/api/register", func(c *fiber.Ctx) error {
		return c.SendString("I'm a POST request!")
	})

	//================================================================

	/*
		Mount Handler
	*/
	/*Gắn micro vào app.
	Khi một yêu cầu được gửi đến app với tiền tố ./john thì
	nó sẽ định tuyến đến micro và được xử lý với micro
	*/
	app.Mount("/john", micro) // GET /john/doe -> 200 OK

	micro.Get("/doe", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	//================================================================

	/*
		MountPath Handler
	*/
	one := fiber.New()
	two := fiber.New()
	three := fiber.New()

	two.Mount("/three", three)
	one.Mount("/two", two)
	app.Mount("/one", one)

	one.MountPath()   // "/one"
	two.MountPath()   // "/one/two"
	three.MountPath() // "/one/two/three"
	app.MountPath()   // ""

	log.Fatal(app.Listen(":3000"))

}
