package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	length := 0
	for i := range str {
		length = i
	}
	fmt.Println(length + 1)
}
