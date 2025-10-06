package main

import (
	"fmt"
	"math"
)

func main() {
	var a float64
	var b float64
	var c float64
	fmt.Scanf("%f %f %f\n", &a, &b, &c)
	fmt.Println(math.Max(math.Max(a, b), c))
}
