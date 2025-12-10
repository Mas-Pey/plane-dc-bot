package bootstrap

import (
	"plane-discord-bot/internal/application"
	"plane-discord-bot/internal/application/api"
)

func RegisterApi() {
	appContainer.RegisterService("planeHandler", &api.PlaneHandler{})
	appContainer.RegisterService("api", &application.Api{})
}
