package main

import (
	"fmt"
)

func genVals() <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))

	go func() {
		defer close(chanStream)

		for i := 0; i < 10; i++ {
			intStream := make(chan interface{}, 1)
			intStream <- i
			close(intStream)
			chanStream <- intStream
		}
	}()

	return chanStream
}

func bridge(done <-chan interface{}, c <-chan (<-chan interface{})) <-chan interface{} {
	bridgeStream := make(chan interface{})

	go func() {
		defer close(bridgeStream)

		for {
			select {
			case <-done:
				return
			case stream, ok := <-c:
				if !ok {
					return
				}

				for val := range stream {
					bridgeStream <- val
				}
			}
		}

	}()

	return bridgeStream
}

func main() {
	channels := genVals()

	for val := range bridge(nil, channels) {
		fmt.Println(val)
	}
}
