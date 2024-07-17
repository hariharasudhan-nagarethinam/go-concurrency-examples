package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randInd(done <-chan interface{}) <-chan interface{} {
	intStream := make(chan interface{})

	go func() {
		defer close(intStream)

		for {
			time.Sleep(1 * time.Second)
			select {
			case <-done:
				return
			case intStream <- rand.Int():
			}
		}
	}()

	return intStream
}

func tee(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for {
			select {
			case <-done:
				return
			case val := <-in:
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
						return
					case out1 <- val:
						fmt.Println("matched out1")
						out1 = nil
					case out2 <- val:
						fmt.Println("matched out2")
						out2 = nil
					default:
						fmt.Println("default")
					}
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	wg := sync.WaitGroup{}
	done := make(chan interface{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(done)
		time.Sleep(2 * time.Second)
	}()

	out1, out2 := tee(done, randInd(done))
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(out1)
		for val := range out1 {
			fmt.Println("out1", val)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for val := range out2 {
			fmt.Println("out2", val)
		}
	}()

	wg.Wait()

}
