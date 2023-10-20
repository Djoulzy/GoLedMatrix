package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	Message       = "Ok"
	StopCharacter = "\r\n\r\n"
)

func handler(conn net.Conn) {

	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			for _, parts := range strings.Fields(data) {
				terminal.AddLine(parts, "#00AA00")
			}
			if isTransportOver(data) {
				break ILOOP
			}
		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}

	}
	w.Write([]byte(Message))
	w.Flush()
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

func Listener() {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(config.HTTPserver.Port))
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", config.HTTPserver.Port, err)
		os.Exit(1)
	}
	defer listen.Close()

	// run loop forever (or until ctrl-c)
	for {
		// get message, output
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handler(conn)
	}
}
