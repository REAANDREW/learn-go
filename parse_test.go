package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

type PacketStringSelector func(packet Packet) string
type StringAssert func(value string, typeId uint16, selector PacketStringSelector)

func assertOnStringPart(t *testing.T) StringAssert {

	return func(value string, typeId uint16, selector PacketStringSelector) {

		buf := new(bytes.Buffer)
		stringValue := value
		partType := typeId
		partLength := uint16(len(stringValue))

		binary.Write(buf, binary.BigEndian, partType)
		binary.Write(buf, binary.BigEndian, partLength)
		buf.WriteString(stringValue)

		packet, err := Parse(buf)
		if err != nil {
			log.Fatalf("err encountered %v", err)
		}
		assert.Equal(t, selector(packet), stringValue, "the string part value does not match the expected")
	}

}

func Test_ParsesTheHostname(t *testing.T) {

	assertOnStringPart(t)("the hostname", 0x0000, func(packet Packet) string {
		return packet.Host.Value
	})

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

	assertOnStringPart(t)("the plugin", 0x0002, func(packet Packet) string {
		return packet.Plugin.Value
	})

}

func Test_ParsesThePluginInstance(t *testing.T) {

	assertOnStringPart(t)("the plugin instance", 0x0003, func(packet Packet) string {
		return packet.PluginInstance.Value
	})

}

func Test_ParsesTheType(t *testing.T) {

	assertOnStringPart(t)("the part type", 0x0004, func(packet Packet) string {
		return packet.Type.Value
	})

}
