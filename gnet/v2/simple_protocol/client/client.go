package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/zlx2019/go-net-v2/simple_protocol/protocol"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9601")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	codec := protocol.SimpleCodec{}

	// Receive server message
	go ReceiveMessage(conn)

	// Send message to server
	std := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := std.ReadLine()
		if string(line) == "exit" {
			return
		}
		content, err := codec.Encode(line)
		if err != nil {
			fmt.Printf("msg encode error: %v \n", err)
		}

		_, err = conn.Write(content)
		if err != nil {
			fmt.Printf("write to server error: %v\n", err)
		}

	}
}

func ReceiveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	welcome, _ := reader.ReadBytes('\n')
	fmt.Printf("server: %s\n", string(welcome))

	for {
		header := make([]byte, 6)
		if _, err := io.ReadFull(reader, header); err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("server is closed.")
			} else {
				fmt.Println("read server msg error:", err)
			}
			return
		}
		magic := binary.BigEndian.Uint16(header[:2])
		if magic != 9501 {
			fmt.Println("magic error:", magic)
			return
		}
		payloadLen := binary.BigEndian.Uint32(header[2:6])
		payload := make([]byte, payloadLen)
		if _, err := io.ReadFull(reader, payload); err != nil {
			fmt.Printf("Error reading payload: %v\n", err)
			return
		}
		fmt.Printf(
			"Received Packet: Magic=%x, Length=%d, Data='%s'\n",
			magic,
			payloadLen,
			string(payload),
		)
	}
}
