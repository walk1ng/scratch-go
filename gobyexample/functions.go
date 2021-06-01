package main

import "fmt"

func plus(a int, b int)int {
	return a+b
}

func main() {
	c := plus(10,20)
	fmt.Println("10+20 =",c)
}