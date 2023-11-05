package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func handler(c *fiber.Ctx) error {
	//Example thui
	return c.Next()
}

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

	//================================================================

	/*
		Group
	*/
	api := app.Group("/api", handler) // /api

	v1 := api.Group("/v1", handler) // /api/v1
	v1.Get("/list", handler)        // /api/v1/list
	v1.Get("/user", handler)        // /api/v1/user

	v2 := api.Group("/v2", handler) // /api/v2
	v2.Get("/list", handler)        // /api/v2/list
	v2.Get("/user", handler)        // /api/v2/user

	//================================================================

	/*
		Server Shutdown
	*/
	//Đợi một thời gian ngắn để máy chủ hoạt động
	time.Sleep(10 * time.Second)
	//Server Shutdown
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server Shutdown: %v", err)
	}

	//ShutdownWithTimeout
	//Tắt máy chủ với thời gian đợi tối đa là 10s
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	//ShuwdownWithContext
	//Tạo context để theo dõi việc đóng máy chủ
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//Theo dõi tín hiệu đóng ứng dụng
	go func() {
		c := make(chan os.Signal, 1) //Kênh c nhận tín hiệu, chứa tối đa 1 giá trị
		/*
			- Đăng ký việc theo dõi tín hiệu đóng
			- Đăng kí theo dõi tín hiệu SIGINT (được gửi khi nhấn Ctrl+C trong terminal)
			và SIGTERM(tín hiệu đóng tổng quán)
		*/
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c // Chờ và nhận tín hiệu từ kênh c
		fmt.Println("Received signal")

		//Hủy bỏ ngữ cảnh để bắt đầu quá trình đóng
		cancel()
	}()
	// Sử dụng ShutdownWithContext để đóng máy chủ
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	fmt.Println("Server has been shut down")

	//================================================================

	/*
		Stack
	*/
	app.Get("/john/:age", handler)
	app.Post("/register", handler)
	data, _ := json.MarshalIndent(app.Stack(), "", " ") //Hàm mã hóa thông tin
	fmt.Println(string(data))
	log.Fatal(app.Listen(":3000"))

}
