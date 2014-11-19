package main

import "fmt"
import "math"

type geometry interface {
	area() float64
	perim() float64
}

type rect struct {
	width,hight float64
}

type cycle struct {
	radius float64
}

func (r rect) area() float64 {
	return r.width * r.hight
}

func (r rect) perim() float64 {
	return 2*r.width + 2*r.hight
}

func (c cycle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c cycle) perim() float64 {
	return 2*c.radius*math.Pi
}

func meaure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	r := rect{10,5}
	c := cycle{5}

	meaure(r)
	meaure(c)
}