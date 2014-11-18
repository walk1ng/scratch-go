package main 

import "fmt"

func sum1(nums ...int) {
	fmt.Print(nums," ")
	total := 0
	for _,num := range(nums) {
		total = total + num
	}
	fmt.Println(total)
}

func sum2(nums ...int)int {
	fmt.Print(nums," ")
	total := 0
	for _,num := range(nums) {
		total = total + num
	}
	return total
}

func main() {
	sum1(1,2)
	sum1(1,2,3)

	nums := []int{1,2,3,4}
	sum1(nums...)

	fmt.Println(sum2(1,2,3,4,5))
	fmt.Println(sum2([]int{1,2,3,4,5,6}...))

}