package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/unusualnorm/echovr_lib/stream"
	"github.com/unusualnorm/echovr_lib/symbols"
)

var SNSConfigSuccessv2Symbol uint64 = symbols.GenerateSymbol("SNSConfigSuccessv2")

type SNSConfigSuccessv2 struct {
	Type   uint64
	ID     uint64
	Config string
}

func (message *SNSConfigSuccessv2) Symbol() uint64 {
	return SNSConfigSuccessv2Symbol
}

func (message *SNSConfigSuccessv2) Deserialize(b []byte) error {
	r := bytes.NewReader(b)

	if err := binary.Read(r, binary.LittleEndian, &message.Type); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &message.ID); err != nil {
		return err
	}

	decompressedBytes, err := stream.ReadZstdCompressedBytes(r)
	if err != nil {
		return err
	}

	decompressedR := bytes.NewReader(decompressedBytes)
	message.Config, err = stream.ReadNullTerminatedString(decompressedR)
	if err != nil {
		return err
	}

	return nil
}

func (message *SNSConfigSuccessv2) Serialize() ([]byte, error) {
	b := bytes.NewBuffer([]byte{})

	if err := binary.Write(b, binary.LittleEndian, message.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(b, binary.LittleEndian, message.ID); err != nil {
		return nil, err
	}

	compressedB := bytes.NewBuffer([]byte{})
	if err := stream.WriteNullTerminatedString(compressedB, message.Config); err != nil {
		return nil, err
	}

	if err := stream.WriteZstdCompressedBytes(b, compressedB.Bytes()); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (message *SNSConfigSuccessv2) String() string {
	return fmt.Sprintf("Config{Type: %v, ID: %v, Config: %v}", message.Type, message.ID, message.Config)
}
