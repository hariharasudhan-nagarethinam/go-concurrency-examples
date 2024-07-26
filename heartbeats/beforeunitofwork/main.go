package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context) (_, __ <-chan interface{}) {
	heartbeatStream := make(chan interface{}, 1)
	resultStream := make(chan interface{})

	go func() {
		defer close(heartbeatStream)
		defer close(resultStream)

		for i := 0; i < 10; i++ {
			select {
			case heartbeatStream <- struct{}{}:
			default:
			}

			select {
			case <-ctx.Done():
				return
			case resultStream <- i:
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

	heartbeatStream, resultStream := doWork(ctx)

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
