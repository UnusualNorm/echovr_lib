package echovr

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/unusualnorm/echovr_lib/stream"
)

var PACKET_HEADER uint64 = 0xBB8CE7A278BB40F6

type Packet struct {
	Header uint64
	Symbol uint64
	Data   []byte
}

func (packet *Packet) Deserialize(b []byte) error {
	r := bytes.NewReader(b)

	if err := binary.Read(r, binary.LittleEndian, &packet.Header); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &packet.Symbol); err != nil {
		return err
	}

	length := uint64(0)
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return err
	}

	err := error(nil)
	packet.Data, err = stream.ReadBytes(r, int(length))
	return err
}

func (packet *Packet) Serialize() ([]byte, error) {
	b := bytes.NewBuffer([]byte{})

	if err := binary.Write(b, binary.LittleEndian, packet.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(b, binary.LittleEndian, packet.Symbol); err != nil {
		return nil, err
	}
	if err := binary.Write(b, binary.LittleEndian, uint64(len(packet.Data))); err != nil {
		return nil, err
	}
	if err := stream.WriteBytes(b, packet.Data); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (packet *Packet) String() string {
	data := ""
	for _, b := range packet.Data {
		data += fmt.Sprintf("%02x", b)
	}
	return fmt.Sprintf("Packet{Header: %v, Symbol: %v, Data: 0x%v}", packet.Header, packet.Symbol, data)
}
