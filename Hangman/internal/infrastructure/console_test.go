package infrastructure

import "testing"

func Test_ConsoleIO_IsCorrectLetter(t *testing.T) {
	io := NewConsoleIO()

	valid := []rune{'а', 'Я', 'ё', '-', ' '}
	invalid := []rune{'A', '1', '$'}

	for _, v := range valid {
		if !io.IsCorrectLetter(v) {
			t.Errorf("expected %q to be valid", v)
		}
	}

	for _, v := range invalid {
		if io.IsCorrectLetter(v) {
			t.Errorf("expected %q to be invalid", v)
		}
	}
}
