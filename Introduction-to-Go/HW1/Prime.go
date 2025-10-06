package main

import "fmt"

func main() {
	var n int
	fmt.Scanf("%d\n", &n)
	for i := 2; i <= n; i++ {
		if IsPrime(i) {
			fmt.Println(i)
		}
	}
}

func IsPrime(x int) bool {
	if x == 1 {
		return false
	} else if x == 2 {
		return true
	} else if x % 2 == 0 {
		return false
	} else {
		for i := 3; i * i <= x; i++ {
			if x % i == 0 {
				return false
			}
		}
	}
	return true
}
