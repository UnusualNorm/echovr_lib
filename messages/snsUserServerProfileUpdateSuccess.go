package messages

import (
	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSUserServerProfileUpdateSuccessSymbol uint64 = echovr.GenerateSymbol("SNSUserServerProfileUpdateSuccess")

type SNSUserServerProfileUpdateSuccess struct {
	UserId echovr.XPlatformID
}

func (m *SNSUserServerProfileUpdateSuccess) Symbol() uint64 {
	return SNSUserServerProfileUpdateSuccessSymbol
}

func (m *SNSUserServerProfileUpdateSuccess) Stream(s *echovr.EasyStream) error {
	return s.StreamStruct(&m.UserId)
}
