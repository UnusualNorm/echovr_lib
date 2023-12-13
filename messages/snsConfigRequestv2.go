package messages

import (
	"fmt"

	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSConfigRequestv2Symbol uint64 = echovr.GenerateSymbol("SNSConfigRequestv2")

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

func (m *SNSConfigRequestv2) Symbol() uint64 {
	return SNSConfigRequestv2Symbol
}

func (m *SNSConfigRequestv2) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{
		func() error { return s.StreamByte(&m.TypeTail) },
		func() error { return s.StreamNullTerminatedJson(&m.ConfigInfo) },
	})
}
