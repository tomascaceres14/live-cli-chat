package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	host = flag.String("host", "localhost", "port")
	port = flag.Int("p", 3090, "port")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))

	if err != nil {
		log.Fatal(err)
	}

	// Canal de control. Tipico para usarlo como semáforo
	channel := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn)
		channel <- struct{}{}
	}()

	CopyContent(conn, os.Stdin)
	conn.Close()

	<-channel
}

func CopyContent(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
