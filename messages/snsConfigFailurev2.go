package messages

import (
	"encoding/binary"
	"fmt"

	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSConfigFailurev2Symbol uint64 = echovr.GenerateSymbol("SNSConfigFailurev2")

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
	Type      uint64
	ID        uint64
	ErrorInfo ConfigErrorInfo
}

func (m *SNSConfigFailurev2) Symbol() uint64 {
	return SNSConfigFailurev2Symbol
}

func (m *SNSConfigFailurev2) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.Type) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.ID) },
		func() error { return s.StreamNullTerminatedJson(&m.ErrorInfo) },
	})
}
