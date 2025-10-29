package infrastructure

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type ConsoleIO struct {
	reader *bufio.Reader
}

func NewConsoleIO() *ConsoleIO {
	return &ConsoleIO{reader: bufio.NewReader(os.Stdin)}
}

func (c *ConsoleIO) Write(text string) {
	fmt.Println(text)
}

func (c *ConsoleIO) ReadLine() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func (c *ConsoleIO) ReadYesOrNo() (string, error) {
	for {
		answer, err := c.ReadLine()
		if err != nil {
			fmt.Println("Input error, please try again.")
			continue
		}

		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer == "yes" || answer == "no" {
			return answer, nil
		}
		fmt.Println("Please enter yes or no!")
	}
}

func (c *ConsoleIO) GetLetter() rune {
	for {
		input, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error, please try again.")
			continue
		}

		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		if utf8.RuneCountInString(input) > 1 {
			fmt.Println("Please enter only one symbol at a time!")
			continue
		}

		letter := []rune(input)[0]
		if !c.IsCorrectLetter(letter) {
			fmt.Println("Please enter russian letter, space or hyphen!")
			continue
		}

		return letter
	}
}

func (c *ConsoleIO) IsCorrectLetter(letter rune) bool {
	return ('а' <= letter && letter <= 'я') || ('А' <= letter && letter <= 'Я') || letter == 'ё' || letter == 'Ё' || letter == ' ' || letter == '-'
}
