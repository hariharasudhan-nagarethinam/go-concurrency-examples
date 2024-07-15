package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
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

func toInt(done <-chan interface{}, value <-chan interface{}) <-chan int {
	intStream := make(chan int)

	go func() {
		defer close(intStream)

		for {
			select {
			case <-done:
				return
			case val := <-value:
				intStream <- val.(int)
			}
		}
	}()

	return intStream
}

func evenNum(done <-chan interface{}, value <-chan int) <-chan int {
	evenStream := make(chan int)

	go func() {
		defer close(evenStream)

		for {
			select {
			case <-done:
				return
			case val := <-value:
				time.Sleep(1 * time.Second)
				if val%2 == 0 {
					evenStream <- val
				}
			}
		}
	}()

	return evenStream
}

func take(done <-chan interface{}, value <-chan int, num int) <-chan int {
	takeStream := make(chan int)

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-value:
			}
		}
	}()

	return takeStream
}

func fanin(done <-chan interface{}, channels []<-chan int) <-chan int {
	mergeStream := make(chan int)
	wg := sync.WaitGroup{}

	mergeFn := func(channel <-chan int) {
		for {
			select {
			case <-done:
				return
			case mergeStream <- <-channel:
			}
		}
	}

	for _, channel := range channels {
		wg.Add(1)
		go mergeFn(channel)
	}

	go func() {
		defer close(mergeStream)
		wg.Wait()
	}()

	return mergeStream

}

func main() {
	done := make(chan interface{})
	defer close(done)

	// fn := func() interface{} { return rand.Intn(500000000) }
	// intStream := toInt(done, repeatFn(done, fn))
	// evenStream := evenNum(done, intStream)
	// takeStream := take(done, evenStream, 10)

	// // before fanout
	// start := time.Now()
	// for val := range takeStream {
	// 	fmt.Println(val)
	// }
	// fmt.Println("Total", time.Since(start))

	// fanout
	fn := func() interface{} { return rand.Intn(500000000) }
	intStream := toInt(done, repeatFn(done, fn))
	numFinders := runtime.NumCPU()
	fanouts := make([]<-chan int, numFinders)
	for i := 0; i < numFinders; i++ {
		fanouts[i] = evenNum(done, intStream)
	}

	faninStream := take(done, fanin(done, fanouts), 10)
	for val := range faninStream {
		fmt.Println(val)
	}
}
