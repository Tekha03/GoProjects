package application

import (
	"hangman/internal/domain"
	"hangman/internal/infrastructure"
	"hangman/pkg/data"
	"testing"
)

func Test_GameService_PickRandomWordFromList(t *testing.T) {
	category := data.FloraFauna
	difficulty := data.Medium
	pair := data.CategoryToDifficulty{Category: category, Difficulty: difficulty}

	words := data.Words[pair]
	if len(words) == 0 {
		t.Fatalf("expected non-empty word list for %v", pair)
	}
}

func Test_GameService_GameStateChangeOnGuess(t *testing.T) {
	h := domain.NewHangman(data.Easy, "", "кот")
	io := infrastructure.NewConsoleIO()
	gs := NewGameService(io, h)

	gs.game.FillLetters('о')
	if h.GetGuessingWord() != "_о_" {
		t.Errorf("expected guessing word to be '_о_', got '%s'", h.GetGuessingWord())
	}

	remaining := h.GetAttemps()
	gs.game.CloserToDeath()
	if h.GetAttemps() != remaining-1 {
		t.Errorf("expected attempts decreased by 1, got %d", h.GetAttemps())
	}
}
