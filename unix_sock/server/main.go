package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
)

// https://go.dev/src/net/unixsock_posix.go
// "unix": SOCK_STREAM
// "unixgram": SOCK_DGRAM
// "unixpacket": SOCK_SEQPACKET

func echoServer(conn net.Conn, pingback bool, ch chan string) {
	for {
		buf := make([]byte, 512)
		nr, err := conn.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		ch <- string(data)

		if pingback {
			// echo back
			fmt.Println("=> server sent")
			_, err := conn.Write([]byte("pong"))
			if err != nil {
				panic("Write: " + err.Error())
			}
		}
	}
}

var opts struct {
	Addr     string `short:"a" long:"addr" default:"/tmp/unix_sock" description:"Address" required:"no"`
	PingBack bool   `short:"p" long:"pingback" description:"Ping back" required:"no"`
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

	l, err := net.Listen("unixpacket", opts.Addr)
	if err != nil {
		fmt.Println("listen error", err.Error())
		return
	}
	defer os.Remove(opts.Addr)

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("accept error", err.Error())
		return
	}
	defer conn.Close()

	go echoServer(conn, opts.PingBack, dataChan)

	for {
		select {
		case <-c:
			fmt.Println("interrupt, cleanup")
			goto END
		case data := <-dataChan:
			fmt.Println("<= server got:", data)
		}
	}

END:
	return
}
