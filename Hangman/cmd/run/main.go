package main

import (
	"fmt"
	"hangman/hangman"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		var player hangman.Hangman
		player.LaunchInteractiveMode()
	} else if len(os.Args) == 3 {
		hidden_word := os.Args[1]
		guessing_word := os.Args[2]

		player := hangman.Constructor("", hidden_word, guessing_word)
		player.NonInteractiveMode()

	} else {
		fmt.Println("Invalid input!")
	}
}
