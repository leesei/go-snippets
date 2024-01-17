package main

import "fmt"

const N int = 217
const BATCH int = 5

func main() {
	c := make(chan int, 6)
	done := make(chan int)

	go process(c, done)
	go produce(c)

	for s := range done {
		println("received", s)
	}
}

func produce(c chan int) {
	for i := 0; i < N; i++ {
		c <- i
	}
	close(c)
}

func process(c chan int, done chan int) {
	count := 0
	sum := 0
	for i := range c {
		// fmt.Printf("[%d] < %d\n", count, i)
		sum += i
		count++
		if count == BATCH {
			fmt.Printf("summed %d items: %d\n", count, sum)
			done <- sum
			count = 0
			sum = 0
		}
	}
	fmt.Printf("summed remaining %d items: %d\n", count, sum)
	done <- sum
	close(done)
}
