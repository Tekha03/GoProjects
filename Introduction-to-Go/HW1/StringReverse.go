package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println(Reversed(str))
}

func Reversed(str string) string {
	str_array := []rune(str)
	length := len(str_array)
	for i, j := length - 1, 0; i > j; i, j = i - 1, j + 1 {
		str_array[i], str_array[j] = str_array[j], str_array[i]
	}
	return string(str_array)
}
