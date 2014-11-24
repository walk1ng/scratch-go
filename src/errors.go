package main 

import "fmt"
import "errors"

func func1(i int) (int, error) {
	if i == 42 {
		return -1, errors.New("can't work with 42")
	} else {
		return i + 3, nil
	}
}

type argError struct {
	arg int
	prob string
}

func (e *argError) Error() string {
	return fmt.Sprintf("%d <-> %s", e.arg, e.prob)
}

func func2(i int) (int, error) {
	if i == 42 {
		return -1, &argError{i, "can't work with it"}
	} else {
		return i + 3, nil
	}
}

func main() {
	for _, val := range([]int{10,42}){
		if r, e := func1(val); e != nil {
			fmt.Println("f1 not worked!", e)
		} else {
			fmt.Println("f1 worked!", r)
		}
	}

	for _, val := range([]int{12,14,42}) {
		if r, e := func2(val); e != nil {
			fmt.Println("f2 not worked!", e)
		} else {
			fmt.Println("f2 worked!", r)
		}
	}

	_, e := func2(42)
	if ae, ok := e.(*argError); ok {
		fmt.Println(ae.arg)
		fmt.Println(ae.prob)
	}
}