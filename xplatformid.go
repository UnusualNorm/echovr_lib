package echovr

import "encoding/binary"

type XPlatformID struct {
	PlatformCode uint64 `json:"platform"`
	AccountID    uint64 `json:"id"`
}

func (m *XPlatformID) Stream(s *EasyStream) error {
	return RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.PlatformCode) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.AccountID) },
	})
}
