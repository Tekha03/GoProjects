package main

import "fmt"

func main() {
	var length int
	fmt.Scan(&length)
	slice := make([]int, length)
	for i := 0; i < length; i++ {
		fmt.Scan(&slice[i])
	}
	var index int
	fmt.Scan(&index)
	fmt.Println(Delete(slice, length, index - 1))
}

func Delete(slice []int, length int, index int) []int {
	copy(slice[index:], slice[index + 1:])
	slice = slice[:length - 1]
	return slice
}
