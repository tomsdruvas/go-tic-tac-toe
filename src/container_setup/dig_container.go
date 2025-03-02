package container_setup

import (
	"go.uber.org/dig"
	"tic-tac-toe-game/src/controllers"
	"tic-tac-toe-game/src/database"
	"tic-tac-toe-game/src/services"
	"tic-tac-toe-game/src/websockets"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	err := container.Provide(controllers.NewHealthCheckController)
	if err != nil {
		return nil
	}

	err = container.Provide(database.NewInMemoryGameSessionDB)
	if err != nil {
		return nil
	}

	err = container.Provide(services.NewGameSessionService)
	if err != nil {
		return nil
	}

	err = container.Provide(controllers.NewGameSessionController)
	if err != nil {
		return nil
	}

	err = container.Provide(websockets.NewWebSocketService)
	if err != nil {
		return nil
	}

	return container
}
