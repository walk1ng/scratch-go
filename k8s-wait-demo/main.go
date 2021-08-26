package main

import (
	"log"

	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	stopper := make(chan struct{})
	wait.Until(worker, time.Second, stopper)
}

func worker() {
	log.Println("start worker")
	time.Sleep(time.Second * 10)
}
