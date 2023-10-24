package main

import (
	"encoding/gob"
	"image"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Djoulzy/GoLedMatrix/clog"
)

const (
	message       = "Ok"
	StopCharacter = "\r\n\r\n"
)

func handler(conn net.Conn) {

	defer conn.Close()

	var (
		buf = Message{}
		r   = gob.NewDecoder(conn)
		w   = gob.NewEncoder(conn)
	)

	mess := Message{
		App:  "TERM",
		Type: "STRING",
		Body: []byte(message),
	}

ILOOP:
	for {
		err := r.Decode(&buf)
		data := string(buf.Type)
		clog.Info("Server", "Handler", "Received: %s", buf.Type)

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			switch buf.Type {
			case "STRING":
				for _, parts := range strings.Fields(buf.Body.(string)) {
					terminal.AddLine(parts, "#00AA00")
					display.Render()
				}
				if isTransportOver(data) {
					break ILOOP
				}
				w.Encode(mess)
			case "IMAGE":
				graphic.SetImage(buf.Body.(image.Image))
				w.Encode(mess)
				display.Render()
			}

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}

	}

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
