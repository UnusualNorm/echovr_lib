package echovr

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/klauspost/compress/zstd"
)

type EasyStream struct {
	mode int // 0 = read, 1 = write
	init []byte
	r    *bytes.Reader
	w    *bytes.Buffer
}

func NewEasyStream(mode int, b []byte) *EasyStream {
	s := &EasyStream{
		mode: mode,
		init: b,
	}

	if mode == 0 {
		s.r = bytes.NewReader(b)
	} else {
		s.w = bytes.NewBuffer(b)
	}

	return s
}

func (s *EasyStream) Bytes() []byte {
	if s.mode == 0 {
		return s.init
	} else {
		return s.w.Bytes()
	}
}

func (s *EasyStream) StreamNumber(order binary.ByteOrder, value any) error {
	if s.mode == 0 {
		return binary.Read(s.r, order, value)
	} else {
		return binary.Write(s.w, order, value)
	}
}

func (s *EasyStream) StreamByte(value *byte) error {
	if s.mode == 0 {
		b, err := s.r.ReadByte()
		if err != nil {
			return err
		}
		*value = b
		return nil
	} else {
		return s.w.WriteByte(*value)
	}
}

func (s *EasyStream) StreamBytes(data []byte, length int) error {
	if s.mode == 0 {
		readBytes := make([]byte, length)
		n, err := s.r.Read(readBytes)
		if err != nil {
			return err
		}
		if n != len(readBytes) {
			return errors.New("StreamBytes: n != len(readBytes)")
		}
		data = readBytes
		return nil
	} else {
		n, err := s.w.Write(data[:length])
		if err != nil {
			return err
		}
		if n != length {
			return errors.New("StreamBytes: n != length")
		}
		return nil
	}
}

func (s *EasyStream) StreamNullTerminatedString(value string) error {
	if s.mode == 0 {
		newValue := ""
		for {
			b, err := s.r.ReadByte()
			if err != nil {
				return err
			}
			if b == 0 {
				break
			}
			newValue += string(b)
		}
		value = newValue
		return nil
	} else {
		valueBytes := []byte(value)
		if err := s.StreamBytes(valueBytes, len(valueBytes)); err != nil {
			return err
		}
		return s.StreamBytes([]byte{0}, 1)
	}
}

func (s *EasyStream) StreamZstdCompressedBytes(data []byte) error {
	if s.mode == 0 {
		decompressedLength := int(0)
		if err := binary.Read(s.r, binary.LittleEndian, &decompressedLength); err != nil {
			return err
		}

		compressedBytes := make([]byte, s.r.Len())
		n, err := s.r.Read(compressedBytes)
		if err != nil {
			return err
		}
		if n != len(compressedBytes) {
			return errors.New("StreamZstdCompressedBytes: n != len(compressedBytes)")
		}

		d, err := zstd.NewReader(nil)
		if err != nil {
			return err
		}

		decompressedBytes, err := d.DecodeAll(compressedBytes, make([]byte, decompressedLength))
		if err != nil {
			return err
		}

		data = decompressedBytes
		return nil
	} else {
		compressedB := bytes.NewBuffer([]byte{})
		if err := s.StreamBytes(data, len(data)); err != nil {
			return err
		}

		if err := s.StreamZstdCompressedBytes(compressedB.Bytes()); err != nil {
			return err
		}

		return nil
	}
}

func (s *EasyStream) StreamZstdEasyStream(stream func(s *EasyStream) error) error {
	if s.mode == 0 {
		decompressedB := []byte{}
		err := s.StreamZstdCompressedBytes(decompressedB)
		if err != nil {
			return err
		}

		return stream(NewEasyStream(0, decompressedB))
	} else {
		compressedB := bytes.NewBuffer([]byte{})
		if err := stream(NewEasyStream(1, compressedB.Bytes())); err != nil {
			return err
		}

		return s.StreamZstdCompressedBytes(compressedB.Bytes())
	}
}

func (s *EasyStream) StreamJson(data interface{}) error {
	if s.mode == 0 {
		jsonString := ""
		err := s.StreamNullTerminatedString(jsonString)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(jsonString), data); err != nil {
			return err
		}

		return nil
	} else {
		jsonString, err := json.Marshal(data)
		if err != nil {
			return err
		}

		if err := s.StreamNullTerminatedString(string(jsonString)); err != nil {
			return err
		}

		return nil
	}
}

func (s *EasyStream) StreamStruct(obj Serializable) error {
	return obj.Stream(s)
}
