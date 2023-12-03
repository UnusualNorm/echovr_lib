package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/unusualnorm/echovr_lib/stream"
	"github.com/unusualnorm/echovr_lib/symbols"
)

var SNSConfigRequestv2Symbol uint64 = symbols.GenerateSymbol("SNSConfigRequestv2")

type ConfigInfo struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (configInfo *ConfigInfo) String() string {
	return fmt.Sprintf("ConfigInfo{Type: \"%v\", ID: \"%v\"}", configInfo.Type, configInfo.ID)
}

func (configInfo *ConfigInfo) Verify() bool {
	return configInfo.Type != "" && configInfo.ID != ""
}

type SNSConfigRequestv2 struct {
	TypeTail   byte
	ConfigInfo ConfigInfo
}

func (message *SNSConfigRequestv2) Symbol() uint64 {
	return SNSConfigRequestv2Symbol
}

func (message *SNSConfigRequestv2) Deserialize(b []byte) error {
	r := bytes.NewReader(b)

	if err := binary.Read(r, binary.LittleEndian, &message.TypeTail); err != nil {
		return err
	}

	configInfoString, err := stream.ReadNullTerminatedString(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(configInfoString), &message.ConfigInfo); err != nil {
		return err
	}

	return nil
}

func (message *SNSConfigRequestv2) Serialize() ([]byte, error) {
	b := bytes.NewBuffer([]byte{})

	if err := binary.Write(b, binary.LittleEndian, message.TypeTail); err != nil {
		return nil, err
	}

	configInfoString, err := json.Marshal(message.ConfigInfo)
	if err != nil {
		return nil, err
	}

	if err := stream.WriteNullTerminatedString(b, string(configInfoString)); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (message *SNSConfigRequestv2) String() string {
	return fmt.Sprintf("SNSConfigRequestv2{TypeTail: 0x%02x, ConfigInfo: %v}", message.TypeTail, message.ConfigInfo.String())
}
