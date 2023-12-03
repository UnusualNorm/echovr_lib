package stream

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/klauspost/compress/zstd"
)

func ReadBytes(r *bytes.Reader, length int) ([]byte, error) {
	value := make([]byte, length)
	n, err := r.Read(value)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, errors.New("ReadBytes: n != length")
	}
	return value, nil
}

func ReadNullTerminatedString(r *bytes.Reader) (string, error) {
	value := ""
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		value += string(b)
	}
	return value, nil
}

func ReadZstdCompressedBytes(r *bytes.Reader) ([]byte, error) {
	decompressedLength := int(0)
	if err := binary.Read(r, binary.LittleEndian, &decompressedLength); err != nil {
		return nil, err
	}

	compressedBytes, err := ReadBytes(r, r.Len())
	if err != nil {
		return nil, err
	}

	d, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}

	decompressedBytes, err := d.DecodeAll(compressedBytes, make([]byte, decompressedLength))
	if err != nil {
		return nil, err
	}

	return decompressedBytes, nil
}

func WriteBytes(b *bytes.Buffer, p []byte) error {
	n, err := b.Write(p)
	if err != nil {
		return err
	}
	if n != len(p) {
		return errors.New("WriteBytes: n != len(p)")
	}
	return nil
}

func WriteNullTerminatedString(b *bytes.Buffer, s string) error {
	n, err := b.WriteString(s)
	if err != nil {
		return err
	}
	if n != len(s) {
		return errors.New("WriteNullTerminatedString: n != len(s)")
	}
	if err := b.WriteByte(0); err != nil {
		return err
	}
	return nil
}

func WriteZstdCompressedBytes(b *bytes.Buffer, p []byte) error {
	e, err := zstd.NewWriter(nil)
	if err != nil {
		return err
	}

	compressedBytes := e.EncodeAll(p, make([]byte, 0, len(p)))

	if err := binary.Write(b, binary.LittleEndian, int32(len(p))); err != nil {
		return err
	}

	if err := WriteBytes(b, compressedBytes); err != nil {
		return err
	}

	return nil
}
