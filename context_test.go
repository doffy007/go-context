package main

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("b"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextA.Value("b"))
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return destination
}

func TestContextWithCancle(t *testing.T) {
	fmt.Println("total goroutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	fmt.Println("total goroutine", runtime.NumGoroutine())
	for v := range destination {
		fmt.Println("counter", v)
		if v == 10 {
			break
		}
	}

	cancel()
	time.Sleep(3 * time.Second)

	fmt.Println("total after", runtime.NumGoroutine())
}

func TestContextWithTimeOut(t *testing.T) {
	fmt.Println("total goroutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("total goroutine", runtime.NumGoroutine())
	for v := range destination {
		fmt.Println("counter", v)
	}

	time.Sleep(3 * time.Second)
	fmt.Println("total after", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("total goroutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("total goroutine", runtime.NumGoroutine())
	for v := range destination {
		fmt.Println("counter", v)
	}

	time.Sleep(3 * time.Second)
	fmt.Println("total after", runtime.NumGoroutine())
}
