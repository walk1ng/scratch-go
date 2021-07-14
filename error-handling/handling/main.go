package main

import (
	"fmt"
	"time"
)

// 一次旅游
type ZooTour1 interface {
	Enter() error                  // 进入
	VisitPanda(panda *Panda) error // 看熊猫
	VisitTiger(tiger *Tiger) error // 看老虎
	Leave() error                  // 离开
}

type MyFunc func(t ZooTour1) error

func NewEnterFunc() MyFunc {
	return func(t ZooTour1) error {
		return t.Enter()
	}
}

func NewVisitPandaFunc(panda *Panda) MyFunc {
	return func(t ZooTour1) error {
		return t.VisitPanda(panda)
	}
}

func NewVisitTigerFunc(tiget *Tiger) MyFunc {
	return func(t ZooTour1) error {
		return t.VisitTiger(tiget)
	}
}

func NewLeaveFunc() MyFunc {
	return func(t ZooTour1) error {
		return t.Leave()
	}
}

type Panda struct{}

type Tiger struct{}

type MyTour struct {
	Name string
}

func (mt MyTour) Enter() error {
	fmt.Printf("%s enter zoo\n", mt.Name)
	time.Sleep(time.Second)
	return nil
}

func (mt MyTour) VisitPanda(panda *Panda) error {
	fmt.Printf("%s visit panda\n", mt.Name)
	time.Sleep(time.Second)
	return fmt.Errorf("panda is missing")
}

func (mt MyTour) VisitTiger(tiger *Tiger) error {
	fmt.Printf("%s visit tiger\n", mt.Name)
	time.Sleep(time.Second)
	return nil

}

func (mt MyTour) Leave() error {
	fmt.Printf("%s leave zoo\n", mt.Name)
	time.Sleep(time.Second)
	return nil
}

func main() {
	mt := MyTour{
		Name: "walk1ng",
	}

	panda := &Panda{}
	tiger := &Tiger{}

	Tour3(mt, panda, tiger)
}

func Tour3(t ZooTour1, panda *Panda, tiger *Tiger) {
	var actions = []MyFunc{
		NewEnterFunc(),
		NewVisitPandaFunc(panda),
		NewVisitTigerFunc(tiger),
		NewLeaveFunc(),
	}

	breakOnError(t, actions)
}

func continueOnError(t ZooTour1, actions []MyFunc) {
	for _, fn := range actions {
		if err := fn(t); err != nil {
			fmt.Printf("%+v\n", err)
			continue
		}
	}
}

func breakOnError(t ZooTour1, actions []MyFunc) {
	for _, fn := range actions {
		if err := fn(t); err != nil {
			fmt.Printf("%+v\n", err)
			break
		}
	}
}
