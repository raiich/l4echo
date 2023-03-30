package payload

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	shaSum  = sha256.Sum256([]byte("l4echo"))
	padding = bytes.Repeat([]byte{'z'}, 1024)
	Size    = len(shaSum) + 8 + 8 + 8 + len(padding)
	endian  = binary.LittleEndian
)

func New() Payload {
	ret := make(Payload, Size)
	copy(ret, shaSum[:])
	copy(ret[Size-len(padding):], padding)
	return ret
}

type Payload []byte

func (p Payload) Seq() uint64 {
	return endian.Uint64(p[len(shaSum):])
}

func (p Payload) Time() time.Time {
	unix := endian.Uint64(p[len(shaSum)+8:])
	nano := endian.Uint64(p[len(shaSum)+8+8:])
	return time.Unix(int64(unix), int64(nano))
}

func (p Payload) SetSeq(seq uint64) {
	endian.PutUint64(p[len(shaSum):], seq)
}

func (p Payload) SetTime(now time.Time) {
	endian.PutUint64(p[len(shaSum)+8:], uint64(now.Unix()))
	endian.PutUint64(p[len(shaSum)+8+8:], uint64(now.Nanosecond()))
}

func Send(w io.Writer, item Payload) error {
	n, err := w.Write(item)
	if err != nil {
		return err
	}
	if n != Size {
		return fmt.Errorf("invalid write size: %d", n)
	}
	return nil
}

func Receive(r io.Reader, buf []byte) (Payload, error) {
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	return Validate(buf[:n])
}

func Validate(data []byte) (Payload, error) {
	if len(data) != Size {
		return nil, fmt.Errorf("invalid payload size: %d", len(data))
	}
	if !bytes.Equal(data[:len(shaSum)], shaSum[:]) {
		return nil, errors.New("invalid payload")
	}
	return data, nil
}
