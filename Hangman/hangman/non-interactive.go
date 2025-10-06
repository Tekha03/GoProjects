package hangman

import "fmt"

func (h *Hangman) NonInteractiveMode() {
	hidden_word_runes := []rune(h.GetHiddenWord())
	guessing_word_runes := []rune(h.GetGuessingWord())

	for i := range hidden_word_runes {
		if hidden_word_runes[i] != guessing_word_runes[i] {
			guessing_word_runes[i] = '*'
		}
	}

	h.SetGuessingWord(guessing_word_runes)
	fmt.Print(h.GetGuessingWord())

	if h.CheckForVictory() {
		fmt.Print(";POS\n")
	} else {
		fmt.Print(";NEG\n")
	}
}
