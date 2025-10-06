package main
import "fmt"

func main() {
	var a int
	var b int
	fmt.Scanf("%d %d\n", &a, &b)
	fmt.Println(add(a, b))
}

func add(a int, b int) int {
	var sum = a + b;
	return sum
}
