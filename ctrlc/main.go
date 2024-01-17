package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-c:
			goto END
		case <-time.After(time.Second):
			fmt.Printf("tick\n")
		}
	}

END:
	fmt.Printf("exit")
}
