package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

const UDP_PACKET_SIZE uint = 65507

func main() {

	uaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", 5555))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", uaddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		packetBytes := make([]byte, UDP_PACKET_SIZE)
		numOfBytesReceived, _, err := conn.ReadFromUDP(packetBytes)
		packetBytes = packetBytes[0:numOfBytesReceived]
		if err != nil {
			log.Fatal(err)
		}
		buffer := bytes.NewBuffer(packetBytes)

		go func(payloadBuffer *bytes.Buffer) {
			packet := Parse(payloadBuffer)
			fmt.Printf("Got a packet %v\n", packet)
		}(buffer)
	}
}
