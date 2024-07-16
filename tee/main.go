package main

import (
	"fmt"
	"math/rand"
	"time"
)

func orDone(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case val, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case valueStream <- val:
				}
			}
		}
	}()

	return valueStream
}

func tee(done <-chan interface{}, c <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, c) {
			var out1, out2 = out1, out2

			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	done := make(chan interface{})
	defer close(done)

	randInt := func() <-chan interface{} {
		randStream := make(chan interface{})

		go func() {
			defer close(randStream)

			time.Sleep(1 * time.Second)
			randStream <- rand.Int()
		}()

		return randStream
	}

	out1, out2 := tee(done, randInt())
	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
