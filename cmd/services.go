package cmd

import (
	"log"
	"plane-discord-bot/internal/bootstrap"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Running services of Plane-Discord-Bot",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(); err != nil {
			log.Println("Env not found")
		} else {
			log.Println("Env succesfully loaded")
		}
		bootstrap.Run()
	},
}

func GetCommand() *cobra.Command {
	return command
}
