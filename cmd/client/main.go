package main

import (
	"encoding/gob"
	"image"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Djoulzy/GoLedMatrix/clog"
	"github.com/Djoulzy/GoLedMatrix/confload"
)

const (
	message1      = "Ping\nPong\nPang"
	message2      = "01_02_03_04_05_06_07_08_09_10"
	StopCharacter = "\r\n\r\n"
)

var config = &confload.ConfigData{}

type Client struct {
	Ip    string
	Port  int
	conn  net.Conn
	read  *gob.Decoder
	write *gob.Encoder
}

func NewClient(ip string, port int) *Client {
	cli := Client{
		Ip:   ip,
		Port: port,
	}
	return &cli
}

func (C *Client) Connect() {
	var err error

	addr := strings.Join([]string{C.Ip, strconv.Itoa(C.Port)}, ":")
	C.conn, err = net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	f, err := os.Open("/Users/jules/go/src/github.com/Djoulzy/GoLedMatrix/media/img/mario.png")
	if err != nil {
		clog.Fatal("scenario", "slideShow", err)
	}
	img, _, _ := image.Decode(f)
	gob.Register(img)

	C.read = gob.NewDecoder(C.conn)
	C.write = gob.NewEncoder(C.conn)
}

func (C *Client) Disconnect() {
	C.conn.Close()
}

func (C *Client) SendImage(file string) {
	f, err := os.Open(file)
	if err != nil {
		clog.Fatal("scenario", "slideShow", err)
	}
	img, _, _ := image.Decode(f)
	gob.Register(img)

	mess := Message{
		App:  "TERM",
		Type: "IMAGE",
		Body: img,
	}
	C.write.Encode(mess)
	buff := make([]byte, 1024)
	C.read.Decode(&buff)
	log.Printf("Receive: %s", buff)
}

func (C *Client) SendText(text string) {
	mess := Message{
		App:  "TERM",
		Type: "STRING",
		Body: text,
	}
	C.write.Encode(mess)
	buff := make([]byte, 1024)
	C.read.Decode(&buff)
	log.Printf("Receive: %s", buff)
}

func main() {
	confload.Load("config.ini", config)
	clog.LogLevel = 5
	clog.StartLogging = true

	cli := NewClient("192.168.0.6", config.HTTPserver.Port)
	cli.Connect()
	defer cli.Disconnect()

	cli.SendImage("/Users/jules/go/src/github.com/Djoulzy/GoLedMatrix/media/img/mario.png")
	cli.SendText(message1)
	cli.SendText(message2)
	cli.SendText(StopCharacter)
}
