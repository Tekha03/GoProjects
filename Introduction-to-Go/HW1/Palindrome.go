package main

import (
	"fmt"
	"bufio"
	"os"
)

func main () {
	str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if IsPalindrome(str) {
		fmt.Println("Palindrome")
	} else {
		fmt.Println("Not Palindrome")
	}
}

func IsPalindrome(str string) bool {
	str_array := []rune(str)
	for i, j := 0, len(str_array) - 2; i < j; i, j = i + 1, j - 1 {
		if str_array[i] != str_array[j] {
			return false
		}
	}
	return true
}
