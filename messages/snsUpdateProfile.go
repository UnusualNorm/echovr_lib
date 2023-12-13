package messages

import (
	"encoding/binary"

	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSUpdateProfileSymbol uint64 = echovr.GenerateSymbol("SNSUpdateProfile")

type SNSUpdateProfile struct {
	Session       echovr.LoginSession
	XPlatformID   echovr.XPlatformID
	Clientprofile string
}

func (m *SNSUpdateProfile) Symbol() uint64 {
	return SNSUpdateProfileSymbol
}

func (m *SNSUpdateProfile) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{
		func() error { return s.StreamStruct(&m.Session) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.XPlatformID) },
		func() error {
			return s.StreamZstdEasyStream(func(decompressedS *echovr.EasyStream) error {
				return decompressedS.StreamNullTerminatedString(&m.Clientprofile)
			})
		},
	})
}
