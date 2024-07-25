package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(
	context context.Context,
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan interface{}) {
	heartbeatStream := make(chan interface{})
	resultStream := make(chan interface{})

	go func() {
		//defer close(heartbeatStream) // simulate panic
		//defer close(resultStream)

		beat := time.Tick(pulseInterval)
		result := time.Tick(pulseInterval * 2)

		for i := 0; i < 2; i++ {
			select {
			case <-context.Done():
				return
			case <-beat:
				select {
				case <-context.Done():
					return
				case heartbeatStream <- struct{}{}:
				default:
				}
			case d := <-result:
				select {
				case <-context.Done():
					return
				case <-beat:
					heartbeatStream <- struct{}{}
				case resultStream <- d:
				}
			}
		}

	}()

	return heartbeatStream, resultStream
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})

	timeout := 2 * time.Second
	heartbeat, results := doWork(ctx, timeout/2)
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
			return
		}
	}
}
