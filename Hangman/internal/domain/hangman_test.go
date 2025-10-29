package domain

import (
	"testing"
	"hangman/pkg/data"
)

func Test_Hangman_CheckLetterCaseInsensitive(t *testing.T) {
	h := NewHangman(data.Easy, "", "кот")

	if !h.CheckLetter('К') {
		t.Errorf("expected CheckLetter to be case-insensitive, got false for 'К'")
	}
}

func Test_Hangman_FillLettersCorrectly(t *testing.T) {
	h := NewHangman(data.Easy, "", "кот")
	h.FillLetters('о')
	got := h.GetGuessingWord()
	want := "_о_"

	if got != want {
		t.Errorf("expected guessing word %s, got %s", want, got)
	}
}

func Test_Hangman_VictoryCheck(t *testing.T) {
	h := NewHangman(data.Easy, "кот", "кот")

	if !h.CheckForVictory() {
		t.Errorf("expected victory to be true, got false")
	}
}

func Test_Hangman_DecrementAttempts(t *testing.T) {
	h := NewHangman(data.Easy, "", "кот")
	initial := h.GetAttemps()
	h.CloserToDeath()
	if h.GetAttemps() != initial-1 {
		t.Errorf("expected attempts %d, got %d", initial-1, h.GetAttemps())
	}
}

func Test_Hangman_GameOverAfterAttempts(t *testing.T) {
	h := NewHangman(data.Easy, "", "кот")

	for i := 0; i < int(data.AttemptsDefault); i++ {
		h.CloserToDeath()
	}

	if h.IsAlive() {
		t.Errorf("expected game to be over after all attempts used")
	}
}

func Test_Hangman_InvalidWordLength(t *testing.T) {
	h := NewHangman(data.Easy, "", "")
	if h.GetHiddenWord() != "" {
		t.Errorf("expected hidden word to be empty for invalid length")
	}
}
