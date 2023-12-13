package messages

import (
	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSUpdateProfileSuccessSymbol uint64 = echovr.GenerateSymbol("SNSUpdateProfileSuccess")

type SNSUpdateProfileSuccess struct {
	UserId echovr.XPlatformID
}

func (m *SNSUpdateProfileSuccess) Symbol() uint64 {
	return SNSUpdateProfileSuccessSymbol
}

func (m *SNSUpdateProfileSuccess) Stream(s *echovr.EasyStream) error {
	return s.StreamStruct(&m.UserId)
}
