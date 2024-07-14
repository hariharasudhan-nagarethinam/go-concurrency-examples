package main

import "fmt"

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	repeatStream := make(chan interface{})

	go func() {
		defer close(repeatStream)
		for {
			for _, val := range values {
				select {
				case <-done:
					return
				case repeatStream <- val:
				}
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

	for val := range take(done, repeat(done, 1), 1) {
		fmt.Println(val)
	}
}
