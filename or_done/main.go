package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randInt() <-chan interface{} {
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)

		for {
			time.Sleep(1 * time.Second)
			valStream <- rand.Int()
		}
	}()

	return valStream
}

func orDone(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case val := <-c:
				valueStream <- val
			}
		}
	}()

	return valueStream
}

func fetchUsers(wg *sync.WaitGroup) {
	defer wg.Done()
	done := make(chan interface{})
	go func() {
		defer close(done)
		time.Sleep(2 * time.Second)
	}()

	for val := range orDone(done, randInt()) {
		fmt.Println("stream1", val)
	}
}

func fetchPosts(wg *sync.WaitGroup) {
	defer wg.Done()
	done := make(chan interface{})
	go func() {
		defer close(done)
		time.Sleep(2 * time.Second)
	}()

	for val := range orDone(done, randInt()) {
		fmt.Println("stream1", val)
	}
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	go fetchPosts(&wg)
	go fetchUsers(&wg)

	wg.Wait()
}
