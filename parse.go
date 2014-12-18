package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

var messageProcessors map[uint16]parser

func init() {
	messageProcessors = createMessageProcessors()
}

func Parse(payloadBuffer *bytes.Buffer) (parsedPacket Packet) {
	var packet Packet

	for payloadBuffer.Len() > 0 {
		partHeader := new(PartHeader)
		binary.Read(payloadBuffer, binary.BigEndian, partHeader)
		partBuffer := bytes.NewBuffer(payloadBuffer.Next(int(partHeader.PartLength) - 4))
		processor, supports := messageProcessors[partHeader.PartType]
		if supports {
			err := processor(&packet, partBuffer)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("%5.d", partHeader.PartType)
		}
	}
	return packet
}

func createMessageProcessors() (processors map[uint16]parser) {
	messageProcessors := make(map[uint16]parser)
	val := []parserGenerator{parseHostname, parseTime, parseHighTime, parsePlugin, parsePluginInstance, parseProcessType, parseProcessTypeInstance, parseInterval, parseValues, parseHighInterval}

	for _, parserGenFunc := range val {
		parserFunc, typeCode := parserGenFunc()
		messageProcessors[typeCode] = parserFunc
	}

	return messageProcessors
}

type parserGenerator func() (parser parser, typeCode uint16)
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

func parseTime() (parser parser, typeCode uint16) {
	code := uint16(0x0001)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		var value int64
		readErr := binary.Read(payload, binary.BigEndian, &value)
		if readErr != nil {
			return readErr
		} else {
			numericPart := NumericPart{partHeaderFromBuffer(code, payload), value}
			packet.Time = numericPart
			return nil
		}
	}, code
}

func parseHighTime() (parser parser, typeCode uint16) {
	code := uint16(0x0008)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		var value int64
		readErr := binary.Read(payload, binary.BigEndian, &value)
		if readErr != nil {
			return readErr
		} else {
			numericPart := NumericPart{partHeaderFromBuffer(code, payload), value >> 30}
			packet.TimeHigh = numericPart
			return nil
		}
	}, code
}

func parsePlugin() (parser parser, typeCode uint16) {
	code := uint16(0x0002)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		stringPart := StringPart{partHeaderFromBuffer(code, payload), payload.String()}
		packet.Plugin = stringPart
		return nil
	}, code
}

func parsePluginInstance() (parser parser, typeCode uint16) {
	code := uint16(0x0003)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		stringPart := StringPart{partHeaderFromBuffer(code, payload), payload.String()}
		packet.PluginInstance = stringPart
		return nil
	}, code
}

func parseProcessType() (parser parser, typeCode uint16) {
	code := uint16(0x0004)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		stringPart := StringPart{partHeaderFromBuffer(code, payload), payload.String()}
		packet.Type = stringPart
		return nil
	}, code
}

func parseProcessTypeInstance() (parser parser, typeCode uint16) {
	code := uint16(0x0005)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		stringPart := StringPart{partHeaderFromBuffer(code, payload), payload.String()}
		packet.TypeInstance = stringPart
		return nil
	}, code
}

func parseInterval() (parser parser, typeCode uint16) {
	code := uint16(0x0008)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		var value int64
		readErr := binary.Read(payload, binary.BigEndian, &value)
		if readErr != nil {
			return readErr
		} else {
			numericPart := NumericPart{partHeaderFromBuffer(code, payload), value}
			packet.Interval = numericPart
			return nil
		}
	}, code
}

func parseValues() (parser parser, typeCode uint16) {
	code := uint16(0x0006)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		header := partHeaderFromBuffer(code, payload)
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
	}, code
}

func parseHighInterval() (parser parser, typeCode uint16) {
	code := uint16(0x0009)
	return func(packet *Packet, payload *bytes.Buffer) (err error) {
		var value int64
		readErr := binary.Read(payload, binary.BigEndian, &value)
		if readErr != nil {
			return readErr
		} else {
			numericPart := NumericPart{partHeaderFromBuffer(0x0009, payload), value >> 30}
			packet.IntervalHigh = numericPart
			return nil
		}
	}, code
}
