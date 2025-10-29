package main

import (
	"fmt"
	"hangman/internal/application"
	"hangman/internal/domain"
	"hangman/internal/infrastructure"
	"hangman/pkg/data"
	"os"
)

func main() {
	io := infrastructure.NewConsoleIO()
	if len(os.Args) == 1 {
		game := domain.NewHangman(data.NotPicked, "", "")
		gs := application.NewGameService(io, game)
		gs.LaunchInteractiveMode()
	} else if len(os.Args) == 3 {
		game := domain.NewHangman(data.NotPicked, os.Args[1], os.Args[2])
		gs := application.NewGameService(io, game)
		gs.NonInteractiveMode()
	} else {
		fmt.Println("Invalid input!")
	}
}
