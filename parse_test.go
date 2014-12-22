package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

type PacketPartAssertion func(packet Packet) bool
type StringPartAssert func(value string, typeId uint16, selector PacketPartAssertion)
type NumericPartAssert func(value int64, typeId uint16, selector PacketPartAssertion)

func assertOnStringPart(t *testing.T) StringPartAssert {

	return func(value string, typeId uint16, selector PacketPartAssertion) {

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
		if !selector(packet) {
			assert.Fail(t, "the string part value does not match the expected")
		}
	}

}

func assertOnNumericPart(t *testing.T) NumericPartAssert {

	return func(value int64, typeId uint16, selector PacketPartAssertion) {

		buf := new(bytes.Buffer)
		partType := typeId
		partLength := uint16(12)

		binary.Write(buf, binary.BigEndian, partType)
		binary.Write(buf, binary.BigEndian, partLength)
		binary.Write(buf, binary.BigEndian, int64(value))

		packet, err := Parse(buf)
		if err != nil {
			log.Fatalf("err encountered %v", err)
		}
		if !selector(packet) {
			assert.Fail(t, "the numeric part value does not match the expected")
		}
	}

}

func Test_ParsesTheHostname(t *testing.T) {
	expected := "the hostname"
	assertOnStringPart(t)(expected, 0x0000, func(packet Packet) bool {
		actual := packet.Host.Value
		return actual == expected
	})

}

func Test_ParsesTheTime(t *testing.T) {
	expected := time.Now().Unix()
	assertOnNumericPart(t)(expected, 0x0001, func(packet Packet) bool {
		actual := packet.Time.Value
		return actual == expected
	})

}

func Test_ParseTheHighDefintiionTime(t *testing.T) {
	expected := time.Now().Unix() << 30
	assertOnNumericPart(t)(expected, 0x0008, func(packet Packet) bool {
		actual := packet.TimeHigh.Value
		return actual == (expected >> 30)
	})

}

func Test_ParsesThePlugin(t *testing.T) {
	expected := "the plugin"
	assertOnStringPart(t)(expected, 0x0002, func(packet Packet) bool {
		actual := packet.Plugin.Value
		return actual == expected
	})

}

func Test_ParsesThePluginInstance(t *testing.T) {
	expected := "the plugin instance"
	assertOnStringPart(t)(expected, 0x0003, func(packet Packet) bool {
		actual := packet.PluginInstance.Value
		return actual == expected
	})

}

func Test_ParsesTheType(t *testing.T) {
	expected := "the part type"
	assertOnStringPart(t)(expected, 0x0004, func(packet Packet) bool {
		actual := packet.Type.Value
		return actual == expected
	})

}
