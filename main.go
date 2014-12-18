package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func createMessageProcessors() (processors map[uint16]parser) {
	messageProcessors := make(map[uint16]parser)
	val := []parserGen{parseHostname, parseTime, parseHighTime, parsePlugin, parsePluginInstance, parseProcessType, parseProcessTypeInstance, parseInterval, parseValues, parseHighInterval}

	for _, parserGenFunc := range val {
		parserFunc, typeCode := parserGenFunc()
		messageProcessors[typeCode] = parserFunc
	}

	return messageProcessors
}

func main() {
	messageProcessors := createMessageProcessors()

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
			packet := new(Packet)
			for payloadBuffer.Len() > 0 {
				partHeader := new(PartHeader)
				binary.Read(payloadBuffer, binary.BigEndian, partHeader)
				partBuffer := bytes.NewBuffer(payloadBuffer.Next(int(partHeader.PartLength) - 4))
				processor, supports := messageProcessors[partHeader.PartType]
				if supports {
					err := processor(packet, partBuffer)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					fmt.Printf("%5.d", partHeader.PartType)
				}
			}
			fmt.Printf("Got a packet ")
			fmt.Print("\n")
		}(buffer)
	}
}
