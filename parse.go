package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type parser func(packet *Packet, payload *bytes.Buffer) (err error)

func plength(length uint16) uint16 {
	return length + 4
}

func partHeaderFromBuffer(partType uint16, payload *bytes.Buffer) PartHeader {
	return PartHeader{partType, plength(uint16(payload.Len()))}
}

func parseHostname() (parser parser, typeCode uint16) {
	code := uint16(0x0000)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		stringPart := StringPart{partHeaderFromBuffer(code, payload), payload.String()}
		packet.Host = stringPart
		return nil
	}, code
}

func parseTime(packet *Packet, payload *bytes.Buffer) (err error) {
	var value int64
	readErr := binary.Read(payload, binary.BigEndian, &value)
	if readErr != nil {
		return readErr
	} else {
		numericPart := NumericPart{partHeaderFromBuffer(0x0001, payload), value}
		packet.Time = numericPart
		return nil
	}
}

func parseHighTime(packet *Packet, payload *bytes.Buffer) (err error) {
	var value int64
	readErr := binary.Read(payload, binary.BigEndian, &value)
	if readErr != nil {
		return readErr
	} else {
		numericPart := NumericPart{partHeaderFromBuffer(0x0008, payload), value >> 30}
		packet.TimeHigh = numericPart
		return nil
	}
}

func parsePlugin(packet *Packet, payload *bytes.Buffer) (err error) {
	stringPart := StringPart{partHeaderFromBuffer(0x0002, payload), payload.String()}
	packet.Plugin = stringPart
	return nil
}

func parsePluginInstance(packet *Packet, payload *bytes.Buffer) (err error) {
	stringPart := StringPart{partHeaderFromBuffer(0x0003, payload), payload.String()}
	packet.PluginInstance = stringPart
	return nil
}

func parseProcessType(packet *Packet, payload *bytes.Buffer) (err error) {
	stringPart := StringPart{partHeaderFromBuffer(0x0004, payload), payload.String()}
	packet.Type = stringPart
	return nil
}

func parseProcessTypeInstance(packet *Packet, payload *bytes.Buffer) (err error) {
	stringPart := StringPart{partHeaderFromBuffer(0x0005, payload), payload.String()}
	packet.TypeInstance = stringPart
	return nil
}

func parseInterval(packet *Packet, payload *bytes.Buffer) (err error) {
	var value int64
	readErr := binary.Read(payload, binary.BigEndian, &value)
	if readErr != nil {
		return readErr
	} else {
		numericPart := NumericPart{partHeaderFromBuffer(0x0008, payload), value}
		packet.Interval = numericPart
		return nil
	}
}

func parseHighInterval(packet *Packet, payload *bytes.Buffer) (err error) {
	var value int64
	readErr := binary.Read(payload, binary.BigEndian, &value)
	if readErr != nil {
		return readErr
	} else {
		numericPart := NumericPart{partHeaderFromBuffer(0x0009, payload), value >> 30}
		packet.IntervalHigh = numericPart
		return nil
	}
}

func parseValues(packet *Packet, payload *bytes.Buffer) (err error) {
	header := partHeaderFromBuffer(0x0006, payload)
	var numberOfValues uint16
	readErr := binary.Read(payload, binary.BigEndian, &numberOfValues)
	if readErr != nil {
		return readErr
	}

	values := make([]Value, numberOfValues)
	counter := uint16(0)
	for counter < numberOfValues {
		var dataType uint8
		readErr := binary.Read(payload, binary.BigEndian, &dataType)
		if readErr != nil {
			return readErr
		}

		switch dataType {
		case 0:
			var value uint32
			readErr := binary.Read(payload, binary.BigEndian, &value)
			if readErr != nil {
				return readErr
			}
			values[counter] = Value{DataType: dataType, CounterValue: value}
		case 1:
			var value float64
			readErr := binary.Read(payload, binary.LittleEndian, &value)
			if readErr != nil {
				return readErr
			}
			values[counter] = Value{DataType: dataType, GaugeValue: value}
		case 2:
			var value int32
			readErr := binary.Read(payload, binary.BigEndian, &value)
			if readErr != nil {
				return readErr
			}
			values[counter] = Value{DataType: dataType, DeriveValue: value}
		case 3:
			var value int32
			readErr := binary.Read(payload, binary.BigEndian, &value)
			if readErr != nil {
				return readErr
			}
			values[counter] = Value{DataType: dataType, AbsoluteValue: value}
		default:
			return fmt.Errorf("unknown value type")
		}
		counter += 1
	}

	packet.Values = ValuePart{header, numberOfValues, values}

	return nil
}
