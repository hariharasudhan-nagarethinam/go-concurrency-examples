package main

import (
	"fmt"
	"time"
)

func doWork(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan interface{}) {
	pulseStream := make(chan interface{})
	resultStream := make(chan interface{})

	go func() {
		defer close(pulseStream)
		defer close(resultStream)

		pulseTick := time.Tick(pulseInterval)
		workTick := time.Tick(pulseInterval * 2)

		for {
			select {
			case <-done:
				return
			case <-pulseTick:
				select {
				case <-done:
					return
				case pulseStream <- struct{}{}:
				default:
				}
			case w := <-workTick:
				select {
				case <-done:
					return
				case <-pulseTick:
					pulseStream <- struct{}{}
				case resultStream <- w:
				}
			}
		}
	}()

	return pulseStream, resultStream
}

func main() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		close(done)
	})

	timeout := 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("Heartbeat")
		case _, ok := <-results:
			if !ok {
				return
			}
			fmt.Println("results")
		case <-time.After(timeout):
			fmt.Println("Timeout")
		}
	}
}
