package messages

import (
	"encoding/binary"

	echovr "github.com/unusualnorm/echovr_lib"
)

type LoggingLevel uint64

const (
	Debug   LoggingLevel = 0x1
	Info    LoggingLevel = 0x2
	Warning LoggingLevel = 0x4
	Error   LoggingLevel = 0x8
	Default LoggingLevel = 0xE
	Any     LoggingLevel = 0xF
)

var SNSRemoteLogSetv3Symbol uint64 = echovr.GenerateSymbol("SNSRemoteLogSetv3")

type SNSRemoteLogSetv3 struct {
	XPlatformID echovr.XPlatformID
	SessionUUID echovr.SessionUUID
	Unk1        uint64
	Unk2        uint64
	LogLevel    LoggingLevel
	Logs        []string
}

func (m *SNSRemoteLogSetv3) Symbol() uint64 {
	return SNSRemoteLogSetv3Symbol
}

func (m *SNSRemoteLogSetv3) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.XPlatformID) },
		func() error { return s.StreamStruct(&m.SessionUUID) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.Unk1) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.Unk2) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.LogLevel) },
		func() error { return s.StreamStringTable(&m.Logs) },
	})
}
