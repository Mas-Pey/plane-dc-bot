package bootstrap

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/runsystemid/gontainer"
)

var appContainer = gontainer.New()

func Run() {
	log.Println("Preparing the system")

	RegisterRest()

	discord := RegisterDiscord()
	defer discord.Close()

	RegisterApi()

	if err := appContainer.Ready(); err != nil {
		log.Fatal("Failed to populate service")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port :%s", port)

	fiberService := appContainer.GetServiceOrNil("fiber").(*fiber.App)

	if err := fiberService.Listen(":" + port); err != nil {
		log.Fatal("Server Shutdown", err)
	}
}
