package echovr

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/klauspost/compress/zstd"
)

type EasyStream struct {
	Mode int // 0 = read, 1 = write
	init []byte
	r    *bytes.Reader
	w    *bytes.Buffer
}

func NewEasyStream(mode int, b []byte) *EasyStream {
	s := &EasyStream{
		Mode: mode,
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
	if s.Mode == 0 {
		return s.init
	} else {
		return s.w.Bytes()
	}
}

func (s *EasyStream) Len() int {
	if s.Mode == 0 {
		return s.r.Len()
	} else {
		return 0
	}
}

func (s *EasyStream) Position() int {
	if s.Mode == 0 {
		return len(s.init) - s.r.Len()
	} else {
		return s.w.Len()
	}
}

func (s *EasyStream) SetPosition(pos int) error {
	if s.Mode == 0 {
		n, err := s.r.Seek(int64(pos), 0)
		if err != nil {
			return err
		}
		if n != int64(pos) {
			return errors.New("SetPosition: n != pos")
		}
		return nil
	} else {
		s.w.Truncate(pos)
		return nil
	}
}

func (s *EasyStream) StreamNumber(order binary.ByteOrder, value any) error {
	if s.Mode == 0 {
		return binary.Read(s.r, order, value)
	} else {
		return binary.Write(s.w, order, value)
	}
}

func (s *EasyStream) StreamByte(value *byte) error {
	if s.Mode == 0 {
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

func (s *EasyStream) StreamBytes(data *[]byte, length int) error {
	if s.Mode == 0 {
		if length == -1 {
			length = s.r.Len()
		}
		readBytes := make([]byte, length)
		n, err := s.r.Read(readBytes)
		if err != nil {
			return err
		}
		if n != len(readBytes) {
			return errors.New("StreamBytes: n != len(readBytes)")
		}
		*data = readBytes
		return nil
	} else {
		n, err := s.w.Write(*data)
		if err != nil {
			return err
		}
		if n != len(*data) {
			return errors.New("StreamBytes: n != len(*data)")
		}
		return nil
	}
}

func (s *EasyStream) StreamString(value *string) error {
	if s.Mode == 0 {
		valueBytes := make([]byte, s.r.Len())
		n, err := s.r.Read(valueBytes)
		if err != nil {
			return err
		}
		if n != len(valueBytes) {
			return errors.New("StreamString: n != len(valueBytes)")
		}
		*value = string(valueBytes)
		return nil
	} else {
		valueBytes := []byte(*value)
		if err := s.StreamBytes(&valueBytes, len(valueBytes)); err != nil {
			return err
		}
		return nil
	}
}

func (s *EasyStream) StreamNullTerminatedString(value *string) error {
	if s.Mode == 0 {
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
		*value = newValue
		return nil
	} else {
		valueBytes := []byte(*value)
		if err := s.StreamBytes(&valueBytes, len(valueBytes)); err != nil {
			return err
		}
		return s.StreamBytes(&[]byte{0}, 1)
	}
}

func (s *EasyStream) StreamZstdCompressedBytes(data *[]byte) error {
	if s.Mode == 0 {
		decompressedLength := uint32(0)
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

		*data = decompressedBytes
		return nil
	} else {
		decompressedB := bytes.NewBuffer([]byte{})
		if err := s.StreamBytes(data, len(*data)); err != nil {
			return err
		}

		decompressedLength := uint32(decompressedB.Len())
		if err := binary.Write(s.w, binary.LittleEndian, &decompressedLength); err != nil {
			return err
		}

		e, err := zstd.NewWriter(nil)
		if err != nil {
			return err
		}

		compressedBytes := e.EncodeAll(decompressedB.Bytes(), make([]byte, int(decompressedLength)))
		if err := s.StreamBytes(&compressedBytes, len(compressedBytes)); err != nil {
			return err
		}

		return nil
	}
}

func (s *EasyStream) StreamZlibCompressedBytes(data *[]byte) error {
	if s.Mode == 0 {
		decompressedLength := uint64(0)
		if err := binary.Read(s.r, binary.LittleEndian, &decompressedLength); err != nil {
			return err
		}

		compressedBytes := make([]byte, s.r.Len())
		n, err := s.r.Read(compressedBytes)
		if err != nil {
			return err
		}
		if n != len(compressedBytes) {
			return errors.New("StreamZlibCompressedBytes: n != len(compressedBytes)")
		}

		d, err := zlib.NewReader(bytes.NewReader(compressedBytes))
		if err != nil {
			return err
		}

		decompressedBytes := make([]byte, decompressedLength)
		n, err = d.Read(decompressedBytes)
		if err != nil {
			return err
		}

		*data = decompressedBytes
		return nil
	} else {
		decompressedB := bytes.NewBuffer([]byte{})
		if err := s.StreamBytes(data, len(*data)); err != nil {
			return err
		}

		decompressedLength := uint64(decompressedB.Len())
		if err := binary.Write(s.w, binary.LittleEndian, &decompressedLength); err != nil {
			return err
		}

		e := zlib.NewWriter(decompressedB)
		if err := e.Close(); err != nil {
			return err
		}

		decompressedBytes := decompressedB.Bytes()
		if err := s.StreamBytes(&decompressedBytes, len(decompressedBytes)); err != nil {
			return err
		}

		return nil
	}
}

func (s *EasyStream) StreamZstdEasyStream(stream func(s *EasyStream) error) error {
	if s.Mode == 0 {
		decompressedB := []byte{}
		err := s.StreamZstdCompressedBytes(&decompressedB)
		if err != nil {
			return err
		}

		return stream(NewEasyStream(0, decompressedB))
	} else {
		decompressedB := bytes.NewBuffer([]byte{})
		if err := stream(NewEasyStream(1, decompressedB.Bytes())); err != nil {
			return err
		}

		decompressedBytes := decompressedB.Bytes()
		return s.StreamZstdCompressedBytes(&decompressedBytes)
	}
}

func (s *EasyStream) StreamZlibEasyStream(stream func(s *EasyStream) error) error {
	if s.Mode == 0 {
		decompressedB := []byte{}
		err := s.StreamZlibCompressedBytes(&decompressedB)
		if err != nil {
			return err
		}

		return stream(NewEasyStream(0, decompressedB))
	} else {
		decompressedB := bytes.NewBuffer([]byte{})
		if err := stream(NewEasyStream(1, decompressedB.Bytes())); err != nil {
			return err
		}

		decompressedBytes := decompressedB.Bytes()
		return s.StreamZlibCompressedBytes(&decompressedBytes)
	}
}

func (s *EasyStream) StreamJson(data interface{}) error {
	if s.Mode == 0 {
		jsonString := ""
		err := s.StreamString(&jsonString)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(jsonString), data); err != nil {
			return err
		}

		return nil
	} else {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		jsonString := string(jsonBytes)
		if err := s.StreamString(&jsonString); err != nil {
			return err
		}

		return nil
	}
}

func (s *EasyStream) StreamNullTerminatedJson(data interface{}) error {
	if s.Mode == 0 {
		jsonString := ""
		err := s.StreamNullTerminatedString(&jsonString)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(jsonString), data); err != nil {
			return err
		}

		return nil
	} else {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		jsonString := string(jsonBytes)
		if err := s.StreamNullTerminatedString(&jsonString); err != nil {
			return err
		}

		return nil
	}
}

func (s *EasyStream) StreamStruct(obj Serializable) error {
	return obj.Stream(s)
}

func (s *EasyStream) StreamStringTable(strings *[]string) error {
	logCount := uint64(len(*strings))
	s.StreamNumber(binary.LittleEndian, &logCount)
	if s.Mode == 0 {
		stringOffsets := make([]uint32, logCount)
		for i := uint64(0); i < logCount; i++ {
			if err := s.StreamNumber(binary.LittleEndian, &stringOffsets[i]); err != nil {
				return err
			}
		}

		stringBytes := []byte{}
		if err := s.StreamBytes(&stringBytes, -1); err != nil {
			return err
		}

		stringS := NewEasyStream(0, stringBytes)
		*strings = []string{}
		for i := uint64(0); i < logCount; i++ {
			offset := stringOffsets[i]
			if err := stringS.SetPosition(int(offset)); err != nil {
				return err
			}

			if err := stringS.StreamNullTerminatedString(&(*strings)[i]); err != nil {
				return err
			}
		}

		return nil
	} else {
		stringS := NewEasyStream(1, []byte{})
		for i := uint64(0); i < logCount; i++ {
			offset := stringS.Position()
			s.StreamNumber(binary.LittleEndian, &offset)
			stringS.StreamNullTerminatedString(&(*strings)[i])
		}

		stringBytes := stringS.Bytes()
		if err := s.StreamBytes(&stringBytes, -1); err != nil {
			return err
		}

		return nil
	}
}
