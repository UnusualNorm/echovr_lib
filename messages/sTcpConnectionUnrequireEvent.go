package messages

import (
	echovr "github.com/unusualnorm/echovr_lib"
)

var STcpConnectionUnrequireEventSymbol uint64 = echovr.GenerateSymbol("STcpConnectionUnrequireEvent")

type STcpConnectionUnrequireEvent struct {
	Unused byte
}

func (m *STcpConnectionUnrequireEvent) Symbol() uint64 {
	return STcpConnectionUnrequireEventSymbol
}

func (m *STcpConnectionUnrequireEvent) Stream(s *echovr.EasyStream) error {
	return s.StreamByte(&m.Unused)
}
