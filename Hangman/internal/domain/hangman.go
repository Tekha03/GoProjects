package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"hangman/pkg/data"
)

type Hangman struct {
	attemptsLeft int32
	difficulty 	  data.Difficulty
	guessingWord string
	hiddenWord   string
}

func NewHangman(difficulty data.Difficulty, guessingWord string, hiddenWord string) *Hangman {
	var GuessingWord string
	if guessingWord != "" {
		GuessingWord = guessingWord
	} else {
		GuessingWord = strings.Repeat("_", utf8.RuneCountInString(hiddenWord))
	}

	return &Hangman{
		attemptsLeft: data.AttemptsDefault,
		difficulty:    difficulty,
		guessingWord: GuessingWord,
		hiddenWord:   hiddenWord,
	}
}

func (h *Hangman) GetAttemps() int32 {
	return h.attemptsLeft
}

func (h *Hangman) DecrementAttempts() {
	h.attemptsLeft -= 1
}

func (h *Hangman) GetHiddenWord() string {
	return h.hiddenWord
}

func (h *Hangman) SetHiddenWord(newHiddenWord string) {
	h.hiddenWord = newHiddenWord
}

func (h *Hangman) GetGuessingWord() string {
	return h.guessingWord
}

func (h *Hangman) SetGuessingWord(newGuessingWord []rune) {
	h.guessingWord = string(newGuessingWord)
}

func (h *Hangman) IsAlive() bool {
	return h.GetAttemps() > 0
}

func (h *Hangman) CheckForVictory() bool {
	return h.GetGuessingWord() == h.GetHiddenWord()
}

func (h *Hangman) CheckLetter(letter rune) bool {
	hidden := h.hiddenWord
	for _, symbol := range hidden {
		if strings.EqualFold(string(symbol), string(letter)) {
			return true
		}
	}
	return false
}

func (h *Hangman) FillLetters(letter rune) {
	hiddenWordRunes := []rune(h.GetHiddenWord())
	guessingWordRunes := []rune(h.GetGuessingWord())

	for i := range hiddenWordRunes {
		if strings.EqualFold(string(hiddenWordRunes[i]), string(letter)) {
			guessingWordRunes[i] = hiddenWordRunes[i]
		}
	}

	h.SetGuessingWord(guessingWordRunes)
}

func (h *Hangman) CloserToDeath() {
	h.DecrementAttempts()
	fmt.Printf("\nYou are one step closer to death. Attempts left: %d\n", h.attemptsLeft)
}
