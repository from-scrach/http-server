package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.fromscratch.sh/http-server/internal"
)

var (
	port = flag.Int("p", 8080, "sets the port")
	host = flag.String("h", "", "sets the host")
)

func handleConnection(conn net.Conn) {

	defer func() {

		if err := conn.Close(); err != nil {
			log.Printf("error: handle: close: %s: %s", conn.RemoteAddr(), err.Error())
		}

	}()

	message, err := internal.ParseHTTPMessage(conn)
	if err != nil {
		log.Printf("error: handle: parse: %s: %s", conn.RemoteAddr(), err.Error())
		return
	}

	fmt.Println(message)

	if _, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")); err != nil {
		log.Printf("error: write: %s: %s", conn.RemoteAddr(), err.Error())
		return
	}

}

func main() {

	listenOn := fmt.Sprintf("%s:%d", *host, *port)
	ctx, cancel := context.WithCancel(context.Background())

	var lc net.ListenConfig
	n, err := lc.Listen(ctx, "tcp", listenOn)
	if err != nil {
		log.Fatalf("error: listen: %s\n", err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		log.Printf("signal: %s received\n", <-c)
		cancel()

		if err := n.Close(); err != nil {
			log.Printf("error: listen: close: %s\n", err.Error())
		}

	}()

	done := false

	for !done {
		select {
		case <-ctx.Done():
			log.Println("termination signal received, exiting listener loop")
			done = true

		default:
			conn, err := n.Accept()
			if err != nil && !errors.Is(err, net.ErrClosed) {
				log.Printf("error: listen-loop: %s\n", err.Error())
				continue
			}

			go handleConnection(conn)
		}
	}

}
