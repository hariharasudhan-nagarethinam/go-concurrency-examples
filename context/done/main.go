package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Wait()
}

func printGreeting(done <-chan interface{}) error {
	msg, err := genGreeting(done)
	if err != nil {
		return err
	}

	fmt.Println(msg)

	return nil
}

func printFarewell(done <-chan interface{}) error {
	msg, err := genFarewall(done)
	if err != nil {
		return err
	}

	fmt.Println(msg)

	return nil
}

func genGreeting(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "Greeting", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func genFarewall(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "Farewell", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("Cancelled")
	case <-time.After(1 * time.Minute):
	}

	return "EN/US", nil
}
