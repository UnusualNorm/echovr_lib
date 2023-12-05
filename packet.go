package echovr

import (
	"encoding/binary"
	"fmt"
)

var PACKET_HEADER uint64 = 0xBB8CE7A278BB40F6

type Packet struct {
	Header uint64
	Symbol uint64
	Data   []byte
}

func (p *Packet) Stream(s *EasyStream) error {
	length := uint64(len(p.Data))

	return RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &p.Header) },
		func() error { return s.StreamNumber(binary.LittleEndian, &p.Symbol) },
		func() error { return s.StreamNumber(binary.LittleEndian, &length) },
		func() error { return s.StreamBytes(p.Data, int(length)) },
	})
}

func (packet *Packet) String() string {
	data := ""
	for _, b := range packet.Data {
		data += fmt.Sprintf("%02x", b)
	}
	return fmt.Sprintf("Packet{Header: %v, Symbol: %v, Data: 0x%v}", packet.Header, packet.Symbol, data)
}
