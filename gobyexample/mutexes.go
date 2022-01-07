package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var readOps uint64
	var writeOps uint64
	var mutex = &sync.Mutex{}

	var state = make(map[int]int)

	// read
	for i := 0; i < 100; i++ {
		go func() {
			total := 0
			for {
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()

				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// write
	for i := 0; i < 10; i++ {
		go func() {
			for {
				key := rand.Intn(5)
				value := rand.Intn(100)
				mutex.Lock()
				state[key] = value
				mutex.Unlock()

				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)
	fmt.Println("read ops:", atomic.LoadUint64(&readOps))
	fmt.Println("write ops:", atomic.LoadUint64(&writeOps))

	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()

}
