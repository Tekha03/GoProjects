package application

import (
	"fmt"
	"hangman/internal/domain"
	"hangman/internal/infrastructure"
	"hangman/pkg/data"
	"math/rand"
)

type GameService struct {
	game *domain.Hangman
	io 	 *infrastructure.ConsoleIO
}

func NewGameService(io *infrastructure.ConsoleIO, game *domain.Hangman) *GameService {
	return &GameService{io: io, game: game}
}

func (gs *GameService) pickDifficulty() data.Difficulty {
	fmt.Println("\nWould you like to pick difficulty? (yes/no)")
	answer, _ := gs.io.ReadYesOrNo()

	if answer == "yes" {
		gs.io.Write("\nPlease pick one of the recommended difficulties: (you should enter the number of the difficulty)")
		gs.io.Write("1. Easy")
		gs.io.Write("2. Medium")
		gs.io.Write("3. Hard")

		picked := gs.getNumber("difficulty")
		switch picked {
		case 1:
			return data.Easy
		case 2:
			return data.Medium
		case 3:
			return data.Hard
		}
	}

	return data.Difficulty(rand.Intn(3))
}

func (gs *GameService) pickCategory() data.Category {
	fmt.Println("\nWould you like to pick category? (yes/no)")
	answer, _ := gs.io.ReadYesOrNo()

	if answer == "yes" {
		fmt.Println("\nPlease pick one of the recommended categories: (you should enter the number of the category)")
		fmt.Println("1. Флора и фауна")
		fmt.Println("2. Города и страны")
		fmt.Println("3. Наука")
		fmt.Println("4. Хобби")
		fmt.Println("5. Поп-культура")

		picked := gs.getNumber("category")
		switch picked {
		case 1:
			return data.FloraFauna
		case 2:
			return data.CitiesCountries
		case 3:
			return data.Science
		case 4:
			return data.Hobby
		case 5:
			return data.PopCulture
		}
	}

	return data.Category(rand.Intn(5))
}

func (gs *GameService) getNumber(catOrDiff string) int {
	var number int

	for {
		_, err := fmt.Scan(&number)
		if err != nil {
			fmt.Println("Input error, please try again.")
			continue
		}

		if gs.CorrectNumber(number, catOrDiff) {
			return number
		}
		fmt.Println("Please enter correct number!")
	}
}

func (gs *GameService) CorrectNumber(number int, catOrDiff string) bool {
	if catOrDiff == "category" {
		if data.MinCategoryNumber <= number && number <= data.MaxCategoryNumber {
			return true
		} else {
			return false
		}
	} else {
		if data.MinDifficulty <= number && number <= data.MaxDifficulty {
			return true
		} else {
			return false
		}
	}
}

func (gs *GameService) LaunchInteractiveMode() {
	category := gs.pickCategory()
	difficulty := gs.pickDifficulty()

	pairCatDiff := data.CategoryToDifficulty{Category: category, Difficulty: difficulty}
	wordsArray := data.Words[pairCatDiff]
	game := domain.NewHangman(difficulty, "", wordsArray[rand.Intn(len(wordsArray))])
	gs.game = game

	fmt.Printf("\nYou picked %s mode in %s category, and now you have %d attempts in total\n",
			difficulty.DifficultyString(), category.CategoryString(), gs.game.GetAttemps())

	gs.launchGuessingIterations()
}

func (gs *GameService) launchGuessingIterations() {
	for gs.game.IsAlive() && !gs.game.CheckForVictory() {
		gs.io.Write(gs.game.GetGuessingWord())
		gs.showHangman()
		gs.io.Write("\nPlease enter the letter:")

		letter := gs.io.GetLetter()
		if gs.game.CheckLetter(letter) {
			gs.game.FillLetters(letter)
		} else {
			gs.game.CloserToDeath()
		}
	}

	if gs.game.CheckForVictory() {
		gs.io.Write("\nCongratulations! You won.")
	} else {
		gs.io.Write("\nYou lost.")
	}
	gs.io.Write(gs.game.GetHiddenWord())
}

func (gs *GameService) showHangman() {
	gs.io.Write(data.Stages[8 - gs.game.GetAttemps()])
}

func (gs *GameService) NonInteractiveMode() {
	hiddenWordRunes := []rune(gs.game.GetHiddenWord())
	guessingWordRunes := []rune(gs.game.GetGuessingWord())

	for i := range hiddenWordRunes {
		if hiddenWordRunes[i] != guessingWordRunes[i] {
			guessingWordRunes[i] = '*'
		}
	}

	gs.game.SetGuessingWord(guessingWordRunes)
	fmt.Print(gs.game.GetGuessingWord())

	if gs.game.CheckForVictory() {
		gs.io.Write(";POS\n")
	} else {
		gs.io.Write(";NEG\n")
	}
}
