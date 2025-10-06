package main

import "fmt"

func main() {
	var length int
	fmt.Scan(&length)
	arr := make([]int, length)
	sum := 0
	for i := 0; i < length; i++ {
		var a int
		fmt.Scan(&a)
		arr[i] = a
		sum += a
	}
	fmt.Println(float32(sum) / float32(length))
}
