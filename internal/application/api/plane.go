package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
)

type IssuePayload struct {
	Identifier      string  `json:"identifier"`
	SequenceId      int     `json:"sequence_id"`
	Name            string  `json:"name"`
	Description     any     `json:"description"`
	DescriptionText *string `json:"description_stripped"`
	Priority        string  `json:"priority"`
	State           struct {
		Name string `json:"name"`
	} `json:"state"`
	Cycle *struct {
		Name string `json:"name"`
	} `json:"cycle"`
	Assignees []struct {
		DisplayName string `json:"display_name"`
	} `json:"assignees"`
}

type PlaneWebhookPayload struct {
	Event  string       `json:"event"`
	Action string       `json:"action"`
	Data   IssuePayload `json:"data"`
}

type PlaneHandler struct {
	DiscordSession *discordgo.Session `inject:"discord"`
}

type PlaneAPI interface {
	HandleWebhook(*fiber.Ctx) error
	Startup() error
	Shutdown() error
}

func (ph *PlaneHandler) Startup() error {
	return nil
}

func (ph *PlaneHandler) Shutdown() error {
	return nil
}

func CheckNull(text string) string {
	if text == "" || text == "None" || text == "null" {
		return "-"
	}
	return text
}

func (ph *PlaneHandler) HandleWebhook(ctx *fiber.Ctx) error {
	rawBody := ctx.Body()

	var payload PlaneWebhookPayload

	if err := json.Unmarshal(rawBody, &payload); err != nil {
		fmt.Println("[ERROR] Failed to parsing JSON from POST trigger: ", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	log.Printf("[DEBUG] Menerima Event: %s | Action: %s\n", payload.Event, payload.Action)

	var err error

	switch payload.Action {
	case "create":
		log.Println("\n[EVENT(CREATE) PAYLOAD]:", string(rawBody))
		err = ph.eventCreateIssue(ctx, payload)
	case "update":
		log.Println("\n[EVENT(UPDATE) PAYLOAD]:", string(rawBody))
		err = ph.eventUpdateIssue(ctx, payload)
	case "delete":
		log.Println("\n[EVENT(DELETE) PAYLOAD]:", string(rawBody))
		err = ph.eventDeleteIssue(ctx, payload)
	default:
		log.Printf("EVENT(WARNING) : %s\n", payload.Action)
		log.Println("[RAW PAYLOAD]:", string(rawBody))
		return ctx.Status(200).JSON(&fiber.Map{
			"status": "ignored",
		})
	}

	if err != nil {
		log.Println("[ERROR] Processing failed", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "Webhook received and processed",
		"action":  payload.Action,
	})
}

func (ph *PlaneHandler) eventCreateIssue(ctx *fiber.Ctx, payload PlaneWebhookPayload) error {
	issueCode := payload.Data.Identifier
	if issueCode == "" {
		projectPrefix := "PLANE"
		issueCode = fmt.Sprintf("%s-%d", projectPrefix, payload.Data.SequenceId)
	}

	cycleName := "-"
	if payload.Data.Cycle != nil {
		cycleName = payload.Data.Cycle.Name
	}

	assigneeNames := "-"
	if len(payload.Data.Assignees) > 0 {
		var names []string
		guildID := "1443521385946157267"
		for _, person := range payload.Data.Assignees {
			mention := ph.findUserMention(guildID, person.DisplayName)
			names = append(names, mention)
		}
		assigneeNames = strings.Join(names, ", ")
	}

	stateName := CheckNull(payload.Data.State.Name)
	priority := CheckNull(payload.Data.Priority)

	description := "-"
	if payload.Data.DescriptionText != nil && *payload.Data.DescriptionText != "" {
		description = *payload.Data.DescriptionText
		if len(description) > 500 {
			description = description[:500] + "..."
		}
	}

	threadTitle := fmt.Sprintf("[%s] %s", issueCode, payload.Data.Name)
	threadContent := fmt.Sprintf(
		"**New Issue Created!**\n\n"+
			"> **Cycle:** %s\n"+
			"> **State:** %s\n"+
			"> **Priority:** %s\n"+
			"> **Assignees:** %s\n\n"+
			"**Description:**\n%s",
		cycleName,
		stateName,
		priority,
		assigneeNames,
		description,
	)

	fmt.Printf("[DEBUG] Creating Thread: %s\n", threadTitle)

	forumChannelID := "1446116820216975380"
	autoArchiveDuration := 1440

	th, err := ph.DiscordSession.ForumThreadStart(
		forumChannelID,
		threadTitle,
		autoArchiveDuration,
		threadContent,
	)

	if err != nil {
		fmt.Printf("[ERROR] Failed to create forum thread: %v\n", err)
		return fmt.Errorf("failed to create forum thread: %w", err)
	}

	fmt.Printf("[SUCCESS] Thread created successfully! Thread ID: %s\n", th.ID)

	return nil
}

func (ph *PlaneHandler) eventUpdateIssue(ctx *fiber.Ctx, payload PlaneWebhookPayload) error {
	return nil
}

func (ph *PlaneHandler) eventDeleteIssue(ctx *fiber.Ctx, payload PlaneWebhookPayload) error {
	return nil
}

func (ph *PlaneHandler) findUserMention(guildID string, planeName string) string {
	targetName := strings.ToLower(planeName)
	members, err := ph.DiscordSession.GuildMembers(guildID, "", 100)
	if err != nil {
		log.Println("[ERROR] Failed to fetch guild members:", err)
		return planeName
	}
	println("[DEBUG] List Members: ", members)

	for _, m := range members {
		if strings.ToLower(m.Nick) == targetName {
			return m.User.Mention()
		}
		if strings.ToLower(m.User.Username) == targetName {
			return m.User.Mention()
		}
		if strings.ToLower(m.User.GlobalName) == targetName {
			return m.User.Mention()
		}
	}
	return planeName
}
