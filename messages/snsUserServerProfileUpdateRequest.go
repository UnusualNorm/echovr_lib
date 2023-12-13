package messages

import (
	"encoding/binary"

	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSUserServerProfileUpdateRequestSymbol uint64 = echovr.GenerateSymbol("SNSUserServerProfileUpdateRequest")

type SNSUserServerProfileUpdateRequest struct {
	XPlatformID echovr.XPlatformID
	UpdateInfo  string
}

func (m *SNSUserServerProfileUpdateRequest) Symbol() uint64 {
	return SNSUserServerProfileUpdateRequestSymbol
}

func (m *SNSUserServerProfileUpdateRequest) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.XPlatformID) },
		func() error {
			return s.StreamZstdEasyStream(func(decompressedS *echovr.EasyStream) error {
				return decompressedS.StreamNullTerminatedJson(&m.UpdateInfo)
			})
		},
	})
}
