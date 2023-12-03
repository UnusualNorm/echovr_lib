package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/unusualnorm/echovr_lib/stream"
	"github.com/unusualnorm/echovr_lib/symbols"
)

var SNSConfigSuccessv2Symbol uint64 = symbols.GenerateSymbol("SNSConfigSuccessv2")

type Config struct {
	Type  string                     `json:"type"`
	ID    string                     `json:"id"`
	Extra map[string]json.RawMessage `json:"-"`
}

func (config *Config) String() string {
	extraString, err := json.Marshal(config.Extra)
	if err != nil {
		return fmt.Sprintf("Config{Type: \"%v\", ID: \"%v\", Extra: ?}", config.Type, config.ID)
	}
	return fmt.Sprintf("Config{Type: \"%v\", ID: \"%v\", Extra: %v}", config.Type, config.ID, extraString)
}

type SNSConfigSuccessv2 struct {
	Type   uint64
	ID     uint64
	Config Config
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
	configString, err := stream.ReadNullTerminatedString(decompressedR)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(configString), &message.Config); err != nil {
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

	configString, err := json.Marshal(message.Config)
	if err != nil {
		return nil, err
	}

	compressedB := bytes.NewBuffer([]byte{})
	if err := stream.WriteNullTerminatedString(compressedB, string(configString)); err != nil {
		return nil, err
	}

	if err := stream.WriteZstdCompressedBytes(b, compressedB.Bytes()); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (message *SNSConfigSuccessv2) String() string {
	return fmt.Sprintf("Config{Type: %v, ID: %v, Config: %v}", message.Type, message.ID, message.Config.String())
}
