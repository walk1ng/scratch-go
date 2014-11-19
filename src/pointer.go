package main

import "fmt"	

func zeroval(ival int) {
	ival = 0
}

func zeroiptr(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("init: ",i)

	zeroval(i)
	fmt.Println("zeroval: ",i)

	zeroiptr(&i)
	fmt.Println("zeroptr: ",i)

	fmt.Println(&i)	
}