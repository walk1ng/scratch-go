package main 

import "fmt"	

type rect struct {
	hight, width int	
}

func (r *rect) area() int {
	return r.hight * r.width
}

func (r rect) perim() int {
	return 2*r.hight + 2*r.width
}

func (r *rect) increas() {
	r.hight += 1
	r.width += 1
}

func main() {
	r := rect{5,10}

	fmt.Println("area: ",r.area())
	fmt.Println("perim: ",r.perim())

	rp := &r
	fmt.Println("area: ",rp.area())
	fmt.Println("perim: ",rp.perim())

	r.increas()
	fmt.Println("area: ",r.area())
	fmt.Println("perim: ",r.perim())

	rp.increas()
	fmt.Println("area: ",r.area())
	fmt.Println("perim: ",r.perim())		
}