package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var atomicInt = &atomic.Int32{}

func main() {
	wg := &sync.WaitGroup{}
	msg := make(chan string)
	wg.Add(2)
	go simulateQueue(wg, msg)

	go func(msg <-chan string) {
		defer wg.Done()

		internalWg := &sync.WaitGroup{}
		stopLight := make(chan struct{}, 5)

		for m := range msg {
			stopLight <- struct{}{}
			internalWg.Add(1)

			message := m
			go func() {
				defer func() {
					internalWg.Done()
					<-stopLight
				}()

				fmt.Println(message)
				fmt.Println("stoplight working", len(stopLight))
				time.Sleep(time.Second * 1)
			}()

		}
		internalWg.Wait()

	}(msg)
	wg.Wait()

}

func simulateQueue(wg *sync.WaitGroup, msg chan<- string) {
	defer func() {
		wg.Done()
		close(msg)
	}()
	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-timer.C:
			for i := 0; i < 100; i++ {
				msg <- fmt.Sprintf("processo %d", atomicInt.Add(1))
			}
		}
	}
}
