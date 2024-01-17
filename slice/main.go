package main

import "fmt"

const N int = 13
const BATCH int = 5

func main() {
	arr := make([]int, N)

	for i := 0; i < BATCH; i++ {
		arr[i] = i
	}
	for i := BATCH; i < N; i++ {
		arr[i] = i
	}

	fmt.Printf("arr: %v\n", arr)
	fmt.Printf("slice: %v\n", arr[0:BATCH])
	fmt.Printf("slice: %v\n", arr[BATCH:N])
}
