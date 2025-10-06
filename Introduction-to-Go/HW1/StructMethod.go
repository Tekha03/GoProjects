package main

import "fmt"

type Rectangle struct{
	width float32
	height float32
}

func (rectange Rectangle) Area() float32 {
	return rectange.width * rectange.height
}

func main() {
	var a float32
	var b float32
	fmt.Scan(&a, &b)
	rectange := Rectangle{a, b}
	fmt.Println(rectange.Area())
}
