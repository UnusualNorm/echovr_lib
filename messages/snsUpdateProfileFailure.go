package messages

import (
	"encoding/binary"

	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSUpdateProfileFailureSymbol uint64 = echovr.GenerateSymbol("SNSUpdateProfileFailure")

type SNSUpdateProfileFailure struct {
	XPlatformID echovr.XPlatformID
	statusCode  uint64 // HTTP Status Code
	Message     string
}

func (m *SNSUpdateProfileFailure) Symbol() uint64 {
	return SNSUpdateProfileFailureSymbol
}

func (m *SNSUpdateProfileFailure) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.XPlatformID) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.statusCode) },
		func() error { return s.StreamNullTerminatedString(&m.Message) },
	})
}
