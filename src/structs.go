package main 

import "fmt"	
type person struct {
	name string
	age int
}

func main() {
	fmt.Println(person{"steve",20})
	fmt.Println(person{name:"alice",age:100})
	fmt.Println(person{name:"joy"})
	fmt.Println(&person{name:"dave",age:18})

	s := person{"sean", 50}
	fmt.Println(s.age)
	s.age = 55
	fmt.Println(s)

	pts := &s
	fmt.Println(*pts)
	fmt.Println(&pts)
	fmt.Println(pts.name)
}