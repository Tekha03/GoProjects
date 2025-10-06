package hangman

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Hangman struct {
	attempts_left int32
	difficulty 	   string
	guessing_word string
	hidden_word   string
}

func Constructor(difficulty string, guessing_word string, hidden_word string) *Hangman {
	var Guessing_word string
	if guessing_word != "" {
		Guessing_word = guessing_word
	} else {
		Guessing_word = strings.Repeat("_", utf8.RuneCountInString(hidden_word))
	}

	return &Hangman{
		attempts_left: 8,
		difficulty: 	difficulty,
		guessing_word: Guessing_word,
		hidden_word:   hidden_word,
	}
}

func (h *Hangman) GetAttemps() int32 {
	return h.attempts_left
}

func (h *Hangman) DecreaseAttempts() {
	h.attempts_left -= 1
}

func (h *Hangman) GetHiddenWord() string {
	return h.hidden_word
}

func (h *Hangman) SetHiddenWord(new_hidden_word string) {
	h.hidden_word = new_hidden_word
}

func (h *Hangman) GetGuessingWord() string {
	return h.guessing_word
}

func (h *Hangman) SetGuessingWord(new_guessing_word []rune) {
	h.guessing_word = string(new_guessing_word)
}

func (h *Hangman) IsAlive() bool {
	return h.GetAttemps() > 0
}

func (h *Hangman) CheckForVictory() bool {
	return h.GetGuessingWord() == h.GetHiddenWord()
}

func (h *Hangman) Win() {
	h.attempts_left = 0
	fmt.Println("\nCongratulations! You won.")
	fmt.Println(h.GetHiddenWord())
}

func (h *Hangman) CloserToDeath() {
	h.DecreaseAttempts()
	if h.IsAlive() {
		fmt.Println("\nYou are one step closer to death.")
	} else {
		h.Death()
	}
}

func (h *Hangman) Death() {
	h.ShowHangman()
	fmt.Println("\nYou lost.")
	fmt.Println(h.GetHiddenWord())
}
