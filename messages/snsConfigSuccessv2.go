package messages

import (
	"encoding/binary"

	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSConfigSuccessv2Symbol uint64 = echovr.GenerateSymbol("SNSConfigSuccessv2")

type SNSConfigSuccessv2 struct {
	Type   uint64
	ID     uint64
	Config string
}

func (m *SNSConfigSuccessv2) Symbol() uint64 {
	return SNSConfigSuccessv2Symbol
}

func (m *SNSConfigSuccessv2) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.Type) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.ID) },
		func() error {
			return s.StreamZstdEasyStream(func(decompressedS *echovr.EasyStream) error {
				return decompressedS.StreamNullTerminatedJson(&m.Config)
			})
		},
	})
}
