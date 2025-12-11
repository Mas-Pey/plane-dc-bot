package bootstrap

import (
	"fmt"
	"log"
	"os"
	"plane-discord-bot/internal/application/commands"

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
	RegisterApplication()

	if err := appContainer.Ready(); err != nil {
		log.Fatal("Failed to populate service")
	}

	cmdHandler := appContainer.GetServiceOrNil("commandHandler").(*commands.CommandHandler)
	cmdHandler.RegisterCompleteCommand()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	printLogo(8080)

	fiberService := appContainer.GetServiceOrNil("fiber").(*fiber.App)

	if err := fiberService.Listen(":" + port); err != nil {
		log.Fatal("Server Shutdown", err)
	}
}

func printLogo(port int) {
	fmt.Printf(`
	 /$$$$$$$  /$$$$$$$$ /$$     /$$
	| $$__  $$| $$_____/|  $$   /$$/
	| $$  \ $$| $$       \  $$ /$$/ 
	| $$$$$$$/| $$$$$     \  $$$$/  
	| $$____/ | $$__/      \  $$/   
	| $$      | $$          | $$    
	| $$      | $$$$$$$$    | $$    
	|__/      |________/    |__/    
	Listening on port : %d`, port)
}
