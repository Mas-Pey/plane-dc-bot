package application

import (
	"plane-discord-bot/internal/application/api"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	App          *fiber.App   `inject:"fiber"`
	PlaneHandler api.PlaneAPI `inject:"planeHandler"`
}

func (a *Api) Startup() error {
	v1 := a.App.Group("api/v1")
	v1.Post("/webhook/plane", a.PlaneHandler.HandleWebhook)
	return nil
}

func (a *Api) Shutdown() error {
	return nil
}
