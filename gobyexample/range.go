package main 
import "fmt"

func main() {
	nums := []int{3,4,5}
	for i,num := range(nums){
		fmt.Println("index:",i,"num:",num)
	}

	kvs := map[string]int{"k1":100,"k2":101}
	for k,v := range(kvs) {
		fmt.Printf("%s -> %d\n", k,v)
	}

	for i,c := range("gog") {
		fmt.Println(i,c)
	}
}