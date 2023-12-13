package echovr

import (
	"encoding/binary"
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
		func() error { return s.StreamBytes(&p.Data, int(length)) },
	})
}

func DeserializePackets(b []byte) ([]Packet, error) {
	s := NewEasyStream(0, b)

	packets := []Packet{}
	for s.Position() < len(b) {
		packet := Packet{}
		if err := packet.Stream(s); err != nil {
			return nil, err
		}
		packets = append(packets, packet)
	}

	return packets, nil
}

func SerializePackets(packets []Packet) ([]byte, error) {
	s := NewEasyStream(1, []byte{})

	for _, packet := range packets {
		if err := packet.Stream(s); err != nil {
			return nil, err
		}
	}

	return s.Bytes(), nil
}

func SerializeMessages(messages []Message) ([]byte, error) {
	packets := []Packet{}

	for _, message := range messages {
		symbol := message.Symbol()
		messageS := NewEasyStream(1, []byte{})
		if err := messageS.StreamStruct(message); err != nil {
			return nil, err
		}

		packet := Packet{
			Header: PACKET_HEADER,
			Symbol: symbol,
			Data:   messageS.Bytes(),
		}
		packets = append(packets, packet)
	}

	return SerializePackets(packets)
}
