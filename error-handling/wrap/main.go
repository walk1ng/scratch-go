package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type MyErrCode int

const (
	ErrorToyNotFoundCode MyErrCode = iota + 1
	ErrorToyHasBeenBorrowedCode
)

var errMap = map[MyErrCode]string{
	ErrorToyNotFoundCode:        "Toy was not found",
	ErrorToyHasBeenBorrowedCode: "Toy has been borrowed",
}

var (
	ErrorToyNotFound        = NewMyError(ErrorToyNotFoundCode)
	ErrorToyHasBeenBorrowed = NewMyError(ErrorToyHasBeenBorrowedCode)
)

type MyError struct {
	Code    MyErrCode
	Message string
}

func (e *MyError) Error() string {
	return e.Message
}

func NewMyError(code MyErrCode) *MyError {
	return &MyError{
		Code:    code,
		Message: errMap[code],
	}
}

func main() {
	Toys := []string{
		"Gozilla",
		"Super Man",
		"Dragon Ball Z",
	}

	for _, ToyName := range Toys {
		fmt.Printf("======= %q start =======\n", ToyName)
		if err := borrowToy(ToyName); err != nil {
			fmt.Printf("%+v\n", err)
		}
		fmt.Printf("======= %q end =======\n", ToyName)
	}
}

func borrowToy(ToyName string) error {
	err := searchToy(ToyName)

	if err != nil {
		var myErr = new(MyError)
		if errors.As(err, &myErr) {
			fmt.Printf("error code is %d, error message is %s\n", myErr.Code, myErr.Message)
		}

		if errors.Is(err, ErrorToyHasBeenBorrowed) {
			fmt.Printf("Toy %q has been borrowed, I will come later!\n", ToyName)
			err = nil
		}
	}

	return err
}

func searchToy(ToyName string) error {
	if len(ToyName) > 10 {
		return errors.Wrapf(ErrorToyNotFound, "Toy name is %s", ToyName)
	} else if len(ToyName) > 8 {
		return errors.WithMessagef(ErrorToyHasBeenBorrowed, "Toy name is %s", ToyName)
	}

	return nil
}
