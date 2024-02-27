package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
)

func reader(r io.Reader, ch chan string) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		ch <- string(buf[0:n])
	}
}

var opts struct {
	Addr string `short:"a" long:"addr" default:"/tmp/unix_sock" description:"Address" required:"no"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}
	fmt.Printf("opts: %+v\n", opts)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	dataChan := make(chan string, 1)

	conn, err := net.Dial("unixpacket", opts.Addr)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	go reader(conn, dataChan)

	for {
		select {
		case <-c:
			fmt.Println("interrupt, cleanup")
			goto END
		case data := <-dataChan:
			fmt.Println("<= client got:", data)
		case <-time.After(1 * time.Second):
			fmt.Println("=> client sent")
			_, err := conn.Write([]byte("ping"))
			if err != nil {
				panic(err.Error())
			}
		}
	}

END:
	return
}
