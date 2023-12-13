package echovr

import "encoding/binary"

type LoginSession struct {
	AccountID uint64 `json:"accountid"`
	Session   uint64 `json:"session"`
}

func (m *LoginSession) Stream(s *EasyStream) error {
	return RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.AccountID) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.Session) },
	})
}

type MatchingSession struct {
	AccountID uint64 `json:"accountid"`
	Session   uint64 `json:"session"`
}

func (m *MatchingSession) Stream(s *EasyStream) error {
	return RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.AccountID) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.Session) },
	})
}

type SessionUUID struct {
	ServerID uint64 `json:"serverid"`
	Session  uint64 `json:"session"`
}

func (m *SessionUUID) Stream(s *EasyStream) error {
	return RunErrorFunctions([]func() error{
		func() error { return s.StreamNumber(binary.LittleEndian, &m.ServerID) },
		func() error { return s.StreamNumber(binary.LittleEndian, &m.Session) },
	})
}
