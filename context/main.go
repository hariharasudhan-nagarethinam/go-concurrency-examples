package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := greetHello(ctx)
		if err != nil {
			fmt.Println("greet hello failed", err)
			cancel()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := greetFarewell(ctx)
		if err != nil {
			fmt.Println("greet farewell failed", err)
			return
		}
	}()

	wg.Wait()
}

func greetHello(ctx context.Context) error {
	msg, err := generateHello(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Greet -->", msg)

	return nil
}

func greetFarewell(ctx context.Context) error {
	msg, err := generateFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Farewell -->", msg)

	return nil
}

func generateHello(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	locale, err := locale(ctx)
	if err != nil {
		return "", err
	}

	switch locale {
	case "EN/US":
		return "greeting", nil
	}

	return "", fmt.Errorf("invalid locale")
}

func generateFarewell(ctx context.Context) (string, error) {
	locale, err := locale(ctx)
	if err != nil {
		return "", err
	}

	switch locale {
	case "EN/US":
		return "farewell", nil
	}

	return "", fmt.Errorf("invalid locale")
}

func locale(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Second * 1):
		return "EN/US", nil
	}
}
