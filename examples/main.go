package main

import (
	grace "github.com/eininst/fiber-grace"
	"github.com/gofiber/fiber/v2"
	"syscall"
	"time"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:     true,
		ReadTimeout: time.Second * 5,
	})

	grace.Listen(app, ":8080", grace.Config{
		Timeout: time.Second * 15,
		Sig:     syscall.SIGILL,
	})
}
