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

func Test_ParseTheHighDefintionTime(t *testing.T) {
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

func Test_ParsesTheTypeInstance(t *testing.T) {
	expected := "the part type instance"
	assertOnStringPart(t)(expected, 0x0005, func(packet Packet) bool {
		actual := packet.TypeInstance.Value
		return actual == expected
	})
}

func Test_ParsesTheValues(t *testing.T) {
	buf := new(bytes.Buffer)
	partType := uint16(0x0006)
	partLength := uint16(16)
	numberOfValues := uint16(2)
	counter := uint16(0)
	binary.Write(buf, binary.BigEndian, partType)
	binary.Write(buf, binary.BigEndian, partLength)
	binary.Write(buf, binary.BigEndian, numberOfValues)
	for counter < numberOfValues {
		counter = counter + 1
		binary.Write(buf, binary.BigEndian, byte(0))
		binary.Write(buf, binary.BigEndian, uint32(5))
	}

	packet, err := Parse(buf)
	if err != nil {
		log.Fatalf("error encountered %v", err)
	}
	assert.Equal(t, len(packet.Values.Values), 2, "incorrect number of values")

	firstValue := packet.Values.Values[0]
	assert.Equal(t, firstValue.DataType, 0, "Data Type mismatch on the ValuePart")
	assert.Equal(t, firstValue.CounterValue, uint32(5), "CounterValue mismatch on the ValuePart")
	secondValue := packet.Values.Values[1]
	assert.Equal(t, secondValue.DataType, 0, "Data Type mismatch on the ValuePart")
	assert.Equal(t, secondValue.CounterValue, uint32(5), "CounterValue mismatch on the ValuePart")
}

func Test_ParsesTheInterval(t *testing.T) {
	expected := int64(1)
	assertOnNumericPart(t)(expected, 0x0007, func(packet Packet) bool {
		actual := packet.Interval.Value
		return actual == expected
	})
}

func Test_ParsesTheHighDefinitionInterval(t *testing.T) {
	expected := int64(1) << 30
	assertOnNumericPart(t)(expected, 0x0009, func(packet Packet) bool {
		actual := packet.IntervalHigh.Value
		return actual == (expected >> 30)
	})
}

func Test_ParsesTheMessage(t *testing.T) {
	expected := "the message"
	assertOnStringPart(t)(expected, 0x0100, func(packet Packet) bool {
		actual := packet.Message.Value
		return actual == expected
	})
}

func Test_ParsesTheSeverity(t *testing.T) {
	expected := int64(1)
	assertOnNumericPart(t)(expected, 0x0101, func(packet Packet) bool {
		actual := packet.Severity.Value
		return actual == expected
	})
}
