package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Djoulzy/GoLedMatrix/clog"
	"github.com/Djoulzy/GoLedMatrix/confload"
)

const (
	message       = "Ping\nPong\nPang"
	StopCharacter = "\r\n\r\n"
)

var config = &confload.ConfigData{}

func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	w.Write([]byte(message))
	w.Flush()
	w.Write([]byte(StopCharacter))
	w.Flush()
	log.Printf("Send: %s", message)

	buff := make([]byte, 1024)
	n, _ := r.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {
	confload.Load("config.ini", config)
	clog.LogLevel = 5
	clog.StartLogging = true

	var (
		ip   = "192.168.0.7"
		port = config.HTTPserver.Port
	)

	SocketClient(ip, port)

}
