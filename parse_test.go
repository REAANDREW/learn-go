package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func Test_ParsesTheHostname(t *testing.T) {

	buf := new(bytes.Buffer)
	hostname := "talula"
	partType := uint16(0x0000)
	partLength := uint16(len(hostname))

	binary.Write(buf, binary.BigEndian, partType)
	binary.Write(buf, binary.BigEndian, partLength)
	buf.WriteString(hostname)

	packet, err := Parse(buf)
	if err != nil {
		log.Fatalf("err encountered %v", err)
	}
	assert.Equal(t, packet.Host.Value, hostname, "the hostname does not match the expected")

}

func Test_ParsesTheTime(t *testing.T) {

	buf := new(bytes.Buffer)
	time := time.Now().Unix()
	partType := uint16(0x0001)
	partLength := uint16(12)

	binary.Write(buf, binary.BigEndian, partType)
	binary.Write(buf, binary.BigEndian, partLength)
	binary.Write(buf, binary.BigEndian, int64(time))

	packet, err := Parse(buf)
	if err != nil {
		log.Fatalf("err encountered %v", err)
	}
	assert.Equal(t, packet.Time.Value, time, "the time does not match the expected")

}

func Test_ParseTheHighDefintiionTIme(t *testing.T) {

	buf := new(bytes.Buffer)
	time := time.Now().Unix() << 30
	partType := uint16(0x0008)
	partLength := uint16(12)

	binary.Write(buf, binary.BigEndian, partType)
	binary.Write(buf, binary.BigEndian, partLength)
	binary.Write(buf, binary.BigEndian, int64(time))

	packet, err := Parse(buf)
	if err != nil {
		log.Fatalf("err encountered %v", err)
	}
	assert.Equal(t, packet.TimeHigh.Value, time>>30, "the time does not match the expected")

}
