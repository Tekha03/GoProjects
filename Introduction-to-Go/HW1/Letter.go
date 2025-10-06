package main

import "fmt"

func main() {
	var symbol string
	fmt.Scanf("%s\n", &symbol)
	if (symbol == "a" || symbol == "e" || symbol == "y" || symbol == "u" || symbol == "i" || symbol == "o" ||
		symbol == "A" || symbol == "E" || symbol == "E" || symbol == "U" || symbol == "I" || symbol == "O") {
		fmt.Println("Vowel")
	} else {
		fmt.Println("Not vowel")
	}
}
