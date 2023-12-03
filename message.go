package echovr

type Message interface {
	Symbol() uint64
	Deserialize([]byte) error
	Serialize() ([]byte, error)
	String() string
}
