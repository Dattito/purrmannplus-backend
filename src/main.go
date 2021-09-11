package main

import (
	"log"

	"github.com/datti-to/purrmannplus-backend/api"
	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/datti-to/purrmannplus-backend/database"
	"github.com/datti-to/purrmannplus-backend/services/signal_message_sender"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	api.InitRoutes(app)

	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	if err := signal_message_sender.Init(); err != nil {
		log.Fatal(err)
	}

	app.Listen(":3000")
}
