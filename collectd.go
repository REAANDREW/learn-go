package main

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
	CounterValue  uint32
	GaugeValue    float64
	DeriveValue   int32
	AbsoluteValue int32
}

type Packet struct {
	Host           StringPart
	Time           NumericPart
	TimeHigh       NumericPart
	Plugin         StringPart
	PluginInstance StringPart
	Type           StringPart
	TypeInstance   StringPart
	Values         ValuePart
	Interval       NumericPart
	IntervalHigh   NumericPart
	Message        StringPart
	Severity       NumericPart
}
