package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

func GenerateRandom(done <-chan int) <-chan int {
	randomInts := make(chan int)

	getRand := func() int {
		num := rand.Intn(500000000)
		return num
	}

	go func() {
		for {
			select {
			case <-done:
				close(randomInts)
				return
			case randomInts <- getRand():
			}
		}
	}()

	return randomInts
}

func PrimeFinder(done <-chan int, stream <-chan int) <-chan int {
	isPrime := func(num int) bool {
		for i := 2; i < num-1; i++ {
			if num%i == 0 {
				return false
			}
		}
		return true
	}

	primeChannel := make(chan int)
	go func() {
		for {
			select {
			case <-done:
				close(primeChannel)
				return
			case number := <-stream:
				if isPrime(number) {
					primeChannel <- number
				}
			}
		}
	}()

	return primeChannel
}

func LimitNums(done <-chan int, stream <-chan int, limit int) <-chan int {
	numsChannel := make(chan int)
	go func() {
		for range limit {
			select {
			case <-done:
				return
			case numsChannel <- <-stream:
			}
		}
		close(numsChannel)
	}()

	return numsChannel
}

func TransferStream(done <-chan int, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	mainChannel := make(chan int)

	forward := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			select {
			case <-done:
				return
			case mainChannel <- v:
			}
		}
	}

	for _, ch := range channels {
		wg.Go(func() {
			forward(ch)
		})
	}

	go func() {
		wg.Wait()
		close(mainChannel)
	}()

	return mainChannel
}

func main() {
	doneChannel := make(chan int)
	defer close(doneChannel)

	randomIntsStream := GenerateRandom(doneChannel)

	cpuCount := runtime.NumCPU()
	fanInStream := make([]<-chan int, cpuCount)
	for i := range cpuCount {
		fanInStream[i] = PrimeFinder(doneChannel, randomIntsStream)
	}

	fanOutStream := TransferStream(doneChannel, fanInStream...)
	limitStream := LimitNums(doneChannel, fanOutStream, 3)

	for i := range limitStream {
		fmt.Println(i)
	}
}
