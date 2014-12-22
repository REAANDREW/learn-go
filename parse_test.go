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

func Test_ParseTheHighDefintiionTime(t *testing.T) {

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
	assert.Equal(t, packet.TimeHigh.Value, time>>30, "the high def time does not match the expected")

}

func Test_ParsesThePlugin(t *testing.T) {

	buf := new(bytes.Buffer)
	plugin := "zeePlugin"
	partType := uint16(0x0002)
	partLength := uint16(len(plugin))

	binary.Write(buf, binary.BigEndian, partType)
	binary.Write(buf, binary.BigEndian, partLength)
	buf.WriteString(plugin)

	packet, err := Parse(buf)
	if err != nil {
		log.Fatalf("err encountered %v", err)
	}
	assert.Equal(t, packet.Plugin.Value, plugin, "the plugin does not match the expected")

}
