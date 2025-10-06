package hangman

import (
	"strings"
	"testing"
)

func TestConstructor(t *testing.T) {
	h := Constructor("easy", "", "кот")

	if h.GetHiddenWord() != "кот" {
		t.Errorf("expected кот, got %s", h.GetHiddenWord())
	}
	if h.GetGuessingWord() != "___" {
		t.Errorf("expected ___, got %s", h.GetGuessingWord())
	}
	if h.GetAttemps() != 8 {
		t.Errorf("expected 8 attempts, got %d", h.GetAttemps())
	}
}

func TestCheckLetter(t *testing.T) {
	h := Constructor("easy", "", "волк")

	if !h.CheckLetter('к') {
		t.Error("letter 'к' must be guessed")
	}
	if h.CheckLetter('з') {
		t.Error("letter 'з' must not be in the word")
	}
}

func TestFillLetter(t *testing.T) {
	h := Constructor("easy", "", "лиса")
	h.FillLetters('и')

	if h.GetGuessingWord() != "_и__" {
		t.Errorf("expected _и__, got %s", h.GetGuessingWord())
	}
}

func TestVictory(t *testing.T) {
	h := Constructor("easy", "", "атом")
	h.FillLetters('а')
	h.FillLetters('т')
	h.FillLetters('о')
	h.FillLetters('м')

	if !h.CheckForVictory() {
		t.Error("player must win, but CheckForVictory returned false")
	}
}

func TestDeath(t *testing.T) {
	h := Constructor("medium", "", "коллекции")

	for i := 0; i < 8; i++ {
		h.CloserToDeath()
	}
	if h.IsAlive() {
		t.Error("player must die after 8 mistakes")
	}
}

func TestCaseInsensitive(t *testing.T) {
	h := Constructor("medium", "", "Хельсинки")

	if !h.CheckLetter('х') {
		t.Error("the game must consider letters without considering their register (Х == х)")
	}
	if !h.CheckLetter('Х') {
		t.Error("the game must consider letters without considering their register (х == Х)")
	}
}

func TestNonInteractiveMode(t *testing.T) {
	h := Constructor("", "к*т", "кот")
	h.SetGuessingWord([]rune("к*т"))

	if !strings.Contains(h.GetGuessingWord(), "*") {
		t.Error("must contain symbol * for unguessed words")
	}
}

func TestGuessingIteration(t *testing.T) {
	h := Constructor("medium", "", "Человек-паук")

	h.GuessingIteration('к')
	if !strings.Contains(h.GetGuessingWord(), "к") {
		t.Error("letter 'к' must have been opened in the word")
	}

	attempts := h.GetAttemps()
	h.GuessingIteration('з')
	if h.GetAttemps() != attempts - 1 {
		t.Error("number of attempts must decrease after mistake")
	}
}

func TestIsCorrectLetter(t *testing.T) {
	h := Constructor("hard", "", "гиппопотам")

	if h.IsCorrectLetter('h') {
		t.Error("only russian letters are allowed")
	}
	if !h.IsCorrectLetter('г') {
		t.Error("letter 'г' must be allowed")
	}
	if h.IsCorrectLetter('1') {
		t.Error("digits must not be considered as correct letters")
	}
	if !h.IsCorrectLetter('-') {
		t.Error("symbol '-' is allowed")
	}
	if !h.IsCorrectLetter(' ') {
		t.Error("space is allowed")
	}
}
