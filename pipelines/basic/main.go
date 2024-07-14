package main

import "fmt"

func Generator(done <-chan interface{}, nums ...int) <-chan int {
	intStream := make(chan int)

	go func() {

		defer close(intStream)

		for _, num := range nums {
			fmt.Println("Gen", num)
			select {
			case <-done:
				return
			case intStream <- num:
			}
		}
	}()

	return intStream
}

func Add(done <-chan interface{}, intStream <-chan int, addtive int) <-chan int {
	addStream := make(chan int)
	go func() {
		defer close(addStream)
		for num := range intStream {
			fmt.Println("Add", num)
			select {
			case <-done:
				return
			case addStream <- num + addtive:
			}
		}
	}()

	return addStream
}

func Mul(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	mulStream := make(chan int)
	go func() {
		defer close(mulStream)
		for num := range intStream {
			fmt.Println("Mul", num)
			select {
			case <-done:
				return
			case mulStream <- num * multiplier:
			}
		}
	}()

	return mulStream
}

func main() {

	done := make(chan interface{})
	defer close(done)

	nums := []int{1, 2, 3, 4}

	intStream := Generator(done, nums...)
	pipeline := Mul(done, Add(done, Mul(done, intStream, 2), 1), 2)

	c := 0
	for data := range pipeline {
		fmt.Println(data)
		if c > 1 {
			done <- true
		}
	}
}
