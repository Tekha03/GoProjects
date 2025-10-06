package main

import "fmt"

func main() {
	var temp float32
	fmt.Scan(&temp)
	fmt.Println(CelsiyToFarengeit(temp))
}

func CelsiyToFarengeit(temp float32) float32 {
	new_temp := temp * 9 / 5 + 32
	return new_temp
}
