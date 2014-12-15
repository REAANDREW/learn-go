package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const UDP_PACKET_SIZE uint = 65507

type PartHeader struct {
	PartType   uint16
	PartLength uint16
}

type StringPart struct {
	Header PartHeader
	Value  string
}

type NumericPart struct {
	Header PartHeader
	Value  int64
}

type ValuePart struct {
	Header         PartHeader
	NumberOfValues uint16
	Values         []Value
}

type Value struct {
	DataType      byte
	CounterValue  uint64
	GaugeValue    float64
	DeriveValue   int64
	AbsoluteValue int64
}

type Packet struct {
	Host           StringPart
	Time           NumericPart
	Plugin         StringPart
	PluginInstance StringPart
	Type           StringPart
	TypeInstance   StringPart
	Values         ValuePart
	Interval       NumericPart
	Message        StringPart
	Severity       NumericPart
}

type part func(packet *Packet, payload *bytes.Buffer) (err error)

func plength(length uint16) uint16 {
	return length + 4
}

func PartHeaderFromBuffer(partType uint16, payload *bytes.Buffer) PartHeader {
	return PartHeader{partType, plength(uint16(payload.Len()))}
}

func hostname(packet *Packet, payload *bytes.Buffer) (err error) {
	stringPart := StringPart{PartHeaderFromBuffer(0x0000, payload), payload.String()}
	packet.Host = stringPart
	log.Printf("type = %d, length = %d, hostname = %s",
		packet.Host.Header.PartType,
		packet.Host.Header.PartLength,
		packet.Host.Value)
	return nil
}

func time(packet *Packet, payload *bytes.Buffer) (err error) {
	var value int64
	readErr := binary.Read(payload, binary.BigEndian, &value)
	if readErr != nil {
		return readErr
	} else {
		numericPart := NumericPart{PartHeaderFromBuffer(0x0001, payload), value}
		packet.Time = numericPart
		log.Printf("type = %d, length = %d, hostname = %s",
			packet.Time.Header.PartType,
			packet.Time.Header.PartLength,
			packet.Time.Value)
		return nil
	}
}

func createMessageProcessors() (processors map[uint16]part) {
	messageProcessors := make(map[uint16]part)
	messageProcessors[0x0000] = hostname
	messageProcessors[0x0001] = time
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

	packet := new(Packet)
	packetBytes := make([]byte, UDP_PACKET_SIZE)

	for {
		numOfBytesReceived, _, err := conn.ReadFromUDP(packetBytes)
		packetBytes = packetBytes[0:numOfBytesReceived]

		if err != nil {
			log.Fatal(err)
		}
		buffer := bytes.NewBuffer(packetBytes)
		go func(payloadBuffer *bytes.Buffer) {
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
					fmt.Print(".")
				}
			}
			fmt.Print("\n")
		}(buffer)
	}
}
