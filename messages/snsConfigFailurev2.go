package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/unusualnorm/echovr_lib/stream"
	"github.com/unusualnorm/echovr_lib/symbols"
)

var SNSConfigFailurev2Symbol uint64 = symbols.GenerateSymbol("SNSConfigFailurev2")

type ConfigErrorInfo struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
	ErrorCode  uint64 `json:"errorCode"`
	Error      string `json:"error"`
}

func (configErrorInfo *ConfigErrorInfo) String() string {
	return fmt.Sprintf("ConfigErrorInfo{Type: \"%v\" Identifier: \"%v\" ErrorCode: 0x%08x Error: \"%v\"}", configErrorInfo.Type, configErrorInfo.Identifier, configErrorInfo.ErrorCode, configErrorInfo.Error)
}

func (configErrorInfo *ConfigErrorInfo) Verify() bool {
	return configErrorInfo.Type != "" && configErrorInfo.Identifier != "" && configErrorInfo.Error != ""
}

type SNSConfigFailurev2 struct {
	Type      uint64 // Unknown
	ID        uint64 // Unknown
	ErrorInfo ConfigErrorInfo
}

func (message *SNSConfigFailurev2) Symbol() uint64 {
	return SNSConfigFailurev2Symbol
}

func (message *SNSConfigFailurev2) Deserialize(b []byte) error {
	r := bytes.NewReader(b)

	if err := binary.Read(r, binary.LittleEndian, &message.Type); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &message.ID); err != nil {
		return err
	}

	errorInfoString, err := stream.ReadNullTerminatedString(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(errorInfoString), &message.ErrorInfo); err != nil {
		return err
	}

	return nil
}

func (message *SNSConfigFailurev2) Serialize() ([]byte, error) {
	b := bytes.NewBuffer([]byte{})

	if err := binary.Write(b, binary.LittleEndian, message.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(b, binary.LittleEndian, message.ID); err != nil {
		return nil, err
	}

	errorInfoString, err := json.Marshal(message.ErrorInfo)
	if err != nil {
		return nil, err
	}

	if err := stream.WriteNullTerminatedString(b, string(errorInfoString)); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (message *SNSConfigFailurev2) String() string {
	return fmt.Sprintf("SNSConfigFailurev2{Type: 0x%08x ID: 0x%08x ErrorInfo: %v}", message.Type, message.ID, message.ErrorInfo.String())
}
