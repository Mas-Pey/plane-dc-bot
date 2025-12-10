package bootstrap

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
)

func RegisterRest() {
	fiberApp := fiber.New()
	appContainer.RegisterService("fiber", fiberApp)
}

func RegisterDiscord() *discordgo.Session {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Token DC Not Found")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to init Discord Bot: ", err)
	}

	// Gateway Configuration
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers

	if err = discord.Open(); err != nil {
		log.Fatal("Failed to connect Discord WebSocket: ", err)
	}

	log.Println("Succesfully connected to Bot Discord")

	appContainer.RegisterService("discord", discord)

	return discord
}
