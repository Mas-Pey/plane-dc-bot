package commands

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	DiscordSession *discordgo.Session `inject:"discord"`
	cmdList        []*discordgo.ApplicationCommand
}

func (ch *CommandHandler) Startup() error {
	log.Println("CommandHandler Ready (Waiting for commands...)")
	return nil
}

func (ch *CommandHandler) Shutdown() error {
	return nil
}

func (ch *CommandHandler) AppendCommand(cmd *discordgo.ApplicationCommand) {
	ch.cmdList = append(ch.cmdList, cmd)
}

func (ch *CommandHandler) RegisterCompleteCommand() {
	log.Println("Registering Slash Commands to Discord...")
	appID := os.Getenv("DISCORD_APP_ID")
	guildID := "1443521385946157267"

	for _, v := range ch.cmdList {
		log.Println("Registering command: ", v.Name)
		_, err := ch.DiscordSession.ApplicationCommandCreate(appID, guildID, v)
		if err != nil {
			log.Printf("Cannot create command %s: %v", v.Name, err)
		}
	}

	log.Println("Slash Commands Registration Done.")
}
