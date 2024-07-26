package main

import (
	"context"
	"fmt"
	"time"
)

func doWorkTimout(
	ctx context.Context,
	pulseInterval time.Duration,
	nums []int,
) (<-chan interface{}, <-chan interface{}) {
	heartbeatStream := make(chan interface{})
	resultStream := make(chan interface{})

	go func() {
		defer close(heartbeatStream)
		defer close(resultStream)

		pulse := time.Tick(pulseInterval)

	numLoop:
		for _, num := range nums {
			fmt.Println(num)
			for {
				select {
				case <-ctx.Done():
					return
				case <-pulse:
					select {
					case heartbeatStream <- struct{}{}:
					default:
					}
				case resultStream <- num:
					continue numLoop
				}
			}
		}
	}()

	return heartbeatStream, resultStream
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})

	heartbeatStream, resultStream := doWorkTimout(ctx, 5*time.Second, []int{1, 2, 3, 4, 5})

	for {
		select {
		case _, ok := <-heartbeatStream:
			if !ok {
				return
			}
			fmt.Println("Received heartbeat")
		case data, ok := <-resultStream:
			if !ok {
				return
			}
			fmt.Println("Received data", data)
		}
	}
}
