package main

import (
	"fmt"
	"math/rand"
)

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	repeatStream := make(chan interface{})

	go func() {
		defer close(repeatStream)
		for {
			select {
			case <-done:
				return
			case repeatStream <- fn():
			}
		}
	}()

	return repeatStream
}

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	fn := func() interface{} {
		return rand.Int()
	}

	for val := range take(done, repeatFn(done, fn), 10) {
		fmt.Println(val)
	}
}
