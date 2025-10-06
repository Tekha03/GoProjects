package hangman

import (
	"bufio"
	"fmt"
	"hangman/data"
	"math/rand"
	"os"
	"strings"
	"unicode/utf8"
)

func (h *Hangman) LaunchInteractiveMode() {
	category := h.PickCategory()
	difficulty := h.PickDifficulty()

	pair_cat_diff := data.CategoryToDifficulty{Category: category, Difficulty: difficulty}
	words_array := data.Words[pair_cat_diff]
	h = Constructor(difficulty, "", words_array[rand.Intn(len(words_array))])

	fmt.Printf("\nYou picked %s mode in %s category, and now you have %d attempts in total\n", difficulty, category, h.GetAttemps())
	h.LaunchGuessingIterations()
}

func (h *Hangman) PickCategory() string {
	fmt.Println("\nWould you like to pick category? (yes/no)")
	answer := h.YesOrNo()

	var category int
	if strings.ToLower(answer) == "yes" {
		fmt.Println("\nPlease pick one of the recommended categories: (you should enter the number of the category)")
		for i := range data.Categories {
			fmt.Printf("%d. %s\n", i + 1, data.Categories[i])
		}

		category = h.GetNumber("category")
	} else if strings.ToLower(answer)== "no" {
		category = rand.Intn(len(data.Categories))
	}

	return data.Categories[category]
}

func (h *Hangman) PickDifficulty() string {
	fmt.Println("\nWould you like to pick difficulty? (yes/no)")
	answer := h.YesOrNo()

	var difficulty int
	if strings.ToLower(answer) == "yes" {
		fmt.Println("\nPlease pick one of the recommended difficulties: (you should enter the number of the difficulty)")
		for i := range data.Difficulties {
			fmt.Printf("%d. %s\n", i + 1, data.Difficulties[i])
		}

		difficulty = h.GetNumber("difficulty")
	} else if strings.ToLower(answer)== "no" {
		difficulty = rand.Intn(len(data.Difficulties))
	}

	return data.Difficulties[difficulty]
}

func (h *Hangman) YesOrNo() string {
	var answer string
	fmt.Scanln(&answer)

	for !h.CorrectInput(answer) {
		fmt.Println("Please enter yes or no!")
		fmt.Scanln(&answer)
	}

	return answer
}

func (h *Hangman) CorrectInput(answer string) bool {
	if strings.ToLower(answer) == "yes" || strings.ToLower(answer) == "no" {
		return true
	} else {
		return false
	}
}

func (h *Hangman) LaunchGuessingIterations() {
	for h.IsAlive() {
		fmt.Println(h.GetGuessingWord())
		h.ShowHangman()
		fmt.Println("\nPlease enter the letter:")

		letter := h.GetLetter()
		h.GuessingIteration(letter)

		if !h.IsAlive() {
			return
		}
	}

	fmt.Println(h.GetHiddenWord())
}

func (h *Hangman) GetNumber(cat_or_diff string) int {
	var number int
	fmt.Scan(&number)

	for !h.CorrectNumber(number, cat_or_diff) {
		fmt.Println("Please enter correct number!")
		fmt.Scan(&number)
	}

	return number - 1
}

func (h *Hangman) GetLetter() rune {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		if utf8.RuneCountInString(input) > 1 {
			fmt.Println("Please enter only one symbol at a time!")
			continue
		}

		letter := []rune(input)[0]
		if !h.IsCorrectLetter(letter) {
			fmt.Println("Please enter russian letter, space or hyphen!")
			continue
		}

		return letter
	}
}

func (h *Hangman) IsCorrectLetter(letter rune) bool {
	return ('а' <= letter && letter <= 'я') || ('А' <= letter && letter <= 'Я') || letter == 'ё' || letter == 'Ё' || letter == ' ' || letter == '-'
}

func (h *Hangman) CorrectNumber(number int, cat_or_diff string) bool {
	if cat_or_diff == "category" {
		if 1 <= number && number <= 5 {
			return true
		} else {
			return false
		}
	} else {
		if 1 <= number && number <= 3 {
			return true
		} else {
			return false
		}
	}
}

func (h *Hangman) GuessingIteration(letter rune) {
	if !h.IsAlive() {
		h.Death()
		return
	} else if h.CheckLetter(letter) {
		h.FillLetters(letter)

		if h.CheckForVictory() {
			h.Win()
			return
		}
	} else {
		h.CloserToDeath()
	}
}

func (h *Hangman) CheckLetter(letter rune) bool {
	hidden := h.GetHiddenWord()
	for _, symbol := range hidden {
		if strings.EqualFold(string(symbol), string(letter)) {
			return true
		}
	}

	return false
}

func (h *Hangman) FillLetters(letter rune) {
	hidden_word_runes := []rune(h.GetHiddenWord())
	guessing_word_runes := []rune(h.GetGuessingWord())

	for i := range hidden_word_runes {
		if strings.EqualFold(string(hidden_word_runes[i]), string(letter)) {
			guessing_word_runes[i] = hidden_word_runes[i]
		}
	}

	h.SetGuessingWord(guessing_word_runes)
}

func (h *Hangman) ShowHangman() {
	fmt.Println(data.Stages[8 - h.attempts_left])
}
