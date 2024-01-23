package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type yamette uint32

func kyaaa(onii chan yamette) {}

func main() {

	app := fiber.New(fiber.Config{
		IdleTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	})

	err := app.Listen("localhost:3265")
	if err != nil {
		panic(err)
	}
}
