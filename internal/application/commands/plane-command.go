package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PlaneCommand struct {
	Handler *CommandHandler `inject:"commandHandler"`
}

func (c *PlaneCommand) Startup() error {

	// root command /plane
	cmd := &discordgo.ApplicationCommand{
		Name:        "plane",
		Description: "Plane Bot slash command",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{

			// SUBCOMMAND 1: ISSUE STATE UPDATE
			// Usage: /plane state [issue_id] [state]
			{
				Name:        "state",
				Description: "Update issue state",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "issue_id",
						Description: "ID Issue (Ex: PLANE-69)",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "state",
						Description: "New state option",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "Backlog", Value: "48590cb9-a347-4c70-977e-3b1cea8ee10d"},
							{Name: "Todo", Value: "c5892faf-2833-4862-a305-afb38c622bb4"},
							{Name: "In Progress", Value: "bc430310-0bc9-4e25-aae4-4d88ad87ede7"},
							{Name: "QA", Value: "37077556-8a4b-48b4-9d08-19a4eb170e23"},
							{Name: "DoneQA", Value: "05c896af-4d36-4919-8822-f298ae2ef8df"},
							{Name: "Done", Value: "244c98a5-f67f-4d70-a546-da7b316e7877"},
							{Name: "Cancelled", Value: "6087999f-b31b-43c0-9186-14ee519a9767"},
						},
					},
				},
			},

			// SUBCOMMAND 2: ASSIGN (skeleton)
			// Usage: /plane assign [issue_id] [user]
			{
				Name:        "assign",
				Description: "Select assignees",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "issue_id",
						Description: "ID Issue (Ex: PLANE-69)",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "user",
						Description: "Assign any",
						Type:        discordgo.ApplicationCommandOptionUser, // User Picker Discord
						Required:    true,
					},
				},
			},
		},
	}

	// Register to Command Handler
	c.Handler.AppendCommand(cmd)

	// Register Handler
	c.Handler.DiscordSession.AddHandler(c.HandleCommand)

	return nil
}

func (c *PlaneCommand) Shutdown() error { return nil }

func (c *PlaneCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == "plane" {
		c.Do(s, i.Interaction)
	}
}

func (c *PlaneCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	options := i.ApplicationCommandData().Options

	subCommand := options[0].Name

	switch subCommand {
	case "state":
		c.handleUpdate(s, i, options[0].Options)
	case "assign":
		c.handleAssign(s, i, options[0].Options)
	default:
		c.respondError(s, i, "Subcommand tidak dikenal")
	}
}

// LOGIC STATE
func (c *PlaneCommand) handleUpdate(s *discordgo.Session, i *discordgo.Interaction, options []*discordgo.ApplicationCommandInteractionDataOption) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	issueID := optionMap["issue_id"].StringValue()
	state := optionMap["state"].StringValue()

	// TODO: Service Plane menyusul lek, alon-alon

	msg := fmt.Sprintf("✅ **Update Berhasil!**\nIssue `%s` sekarang statusnya **%s**", issueID, state)
	c.respondSuccess(s, i, msg)
}

// LOGIC ASSIGN (Simulasi)
func (c *PlaneCommand) handleAssign(s *discordgo.Session, i *discordgo.Interaction, options []*discordgo.ApplicationCommandInteractionDataOption) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	issueID := optionMap["issue_id"].StringValue()
	targetUser := optionMap["user"].UserValue(s)

	// TODO: Service Plane menyusul lek, alon-alon

	msg := fmt.Sprintf("✅ **Assign Berhasil!**\nIssue `%s` diserahkan ke **%s**", issueID, targetUser.Mention())
	c.respondSuccess(s, i, msg)
}

func (c *PlaneCommand) respondSuccess(s *discordgo.Session, i *discordgo.Interaction, content string) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func (c *PlaneCommand) respondError(s *discordgo.Session, i *discordgo.Interaction, content string) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "❌ " + content},
	})
}
