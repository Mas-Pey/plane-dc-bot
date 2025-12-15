package bootstrap

import (
	"plane-discord-bot/internal/application"
	"plane-discord-bot/internal/application/api"
	"plane-discord-bot/internal/application/commands"
	"plane-discord-bot/internal/application/services"
)

func RegisterApi() {
	appContainer.RegisterService("planeHandler", &api.PlaneHandler{})
	appContainer.RegisterService("api", &application.Api{})
}

func RegisterApplication() {
	appContainer.RegisterService("commandHandler", &commands.CommandHandler{})
	appContainer.RegisterService("planeCommand", &commands.PlaneCommand{})

	appContainer.RegisterService("planeService", services.NewPlaneService())
}
