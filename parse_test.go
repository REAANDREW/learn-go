package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_ParsesTheHostname(t *testing.T) {

	buf := new(bytes.Buffer)
	hostname := "talula"
	partType := uint16(0)
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
