package messages

import (
	"fmt"

	"github.com/unusualnorm/echovr_lib/symbols"
)

var STcpConnectionUnrequireEventSymbol uint64 = symbols.GenerateSymbol("STcpConnectionUnrequireEvent")

type STcpConnectionUnrequireEvent struct {
	Unused byte
}

func (message *STcpConnectionUnrequireEvent) Symbol() uint64 {
	return STcpConnectionUnrequireEventSymbol
}

func (message *STcpConnectionUnrequireEvent) Deserialize(b []byte) error {
	if len(b) < 1 {
		return fmt.Errorf("STcpConnectionUnrequireEvent: len(b) < 1")
	}

	message.Unused = b[0]
	return nil
}

func (message *STcpConnectionUnrequireEvent) Serialize() ([]byte, error) {
	return []byte{message.Unused}, nil
}

func (message *STcpConnectionUnrequireEvent) String() string {
	return fmt.Sprintf("STcpConnectionUnrequireEvent{Unused: 0x%02x}", message.Unused)
}
