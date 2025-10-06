package main
import "fmt"

func main() {
	var x int
	fmt.Scanf("%d\n", &x)
	EvenOdd(x)
}

func EvenOdd(x int) {
	if x % 2 == 0 {
		fmt.Println("Even")
	} else {
		fmt.Println("Odd")
	}
}
