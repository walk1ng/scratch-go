package main

import (
	"errors"
	"fmt"
	"go-transaction/pkg/exception"
)

func main() {
	err := demo()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("great!")
}

func demo() (err error) {
	exception.Block{
		Try: func() {
			beginTransaction()
			if err := one(); err != nil {
				panic(err)
			}
			if err := two(); err != nil {
				panic(err)
			}
			if err := three(); err != nil {
				panic(err)
			}
			if err := four(); err != nil {
				panic(err)
			}
			if err := five(); err != nil {
				panic(err)
			}
			commit()
			err = nil
		},
		Catch: func(e interface{}) {
			rollback()
			fmt.Printf("%v panic\n", e)
			err = fmt.Errorf("%v", e)
		},
		Finally: func() {
		},
	}.Do()
	return err
}

//开启事务
func beginTransaction() {
	fmt.Println("beginTransaction")
}

//回滚事务
func rollback() {
	fmt.Println("rollback")
}

//提交事务
func commit() {
	fmt.Println("commit")
}

//执行one操作
func one() (err error) {
	fmt.Println("one ok")
	return nil
}

//执行two操作
func two() (err error) {
	fmt.Println("two ok")
	return nil
}

//执行three操作
func three() (err error) {
	fmt.Println("three ok")
	return nil
}

//执行four操作
func four() (err error) {
	fmt.Println("four ok")
	return nil
}

//执行five操作
func five() (err error) {
	err = errors.New("five panic")
	// panic("five")
	return err
}
