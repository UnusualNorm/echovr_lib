package messages

import (
	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSChannelInfoRequestSymbol uint64 = echovr.GenerateSymbol("SNSChannelInfoRequest")

type SNSChannelInfoRequest struct {
	Unused byte
}

func (m *SNSChannelInfoRequest) Symbol() uint64 {
	return SNSChannelInfoRequestSymbol
}

func (m *SNSChannelInfoRequest) Stream(s *echovr.EasyStream) error {
	return s.StreamByte(&m.Unused)
}
