package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"io"
)

var ErrMagicMismatch = errors.New("magic mismatch")
var ErrIncompletePacket = errors.New("incomplete packet")

var magicBytes []byte

const (
	magicValue       = 9501
	magicBytesSize   = 2
	payloadBytesSize = 4
	payloadOffset    = magicBytesSize + payloadBytesSize
)

func init() {
	magicBytes = make([]byte, magicBytesSize)
	binary.BigEndian.PutUint16(magicBytes, uint16(magicValue))
}

// Protocol format:
//
// * 0           2                       6
// * +-----------+-----------------------+
// * |   magic   |       payload len     |
// * +-----------+-----------+-----------+
// * |                                   |
// * +                                   +
// * |           payload bytes           |
// * +                                   +
// * |            ... ...                |
// * +-----------------------------------+

type SimpleCodec struct{}

// Encode 数据封包
func (codec SimpleCodec) Encode(buf []byte) ([]byte, error) {
	totalSize := payloadOffset + len(buf)

	msg := make([]byte, totalSize)
	copy(msg, magicBytes)

	binary.BigEndian.PutUint32(msg[magicBytesSize:payloadOffset], uint32(len(buf)))
	copy(msg[payloadOffset:totalSize], buf)
	return msg, nil
}

// Unpack 从输入流中解包
func (codec SimpleCodec) Unpack(buf []byte) ([]byte, error) {
	if len(buf) < payloadOffset {
		return nil, ErrIncompletePacket
	}
	if !bytes.Equal(magicBytes, buf[:magicBytesSize]) {
		return nil, ErrMagicMismatch
	}
	payloadSize := binary.BigEndian.Uint32(buf[magicBytesSize:payloadOffset])
	totalSize := payloadOffset + int(payloadSize)
	if len(buf) < totalSize {
		return nil, ErrIncompletePacket
	}
	return buf[payloadOffset:totalSize], nil
}

// Decode 解码
func (codec SimpleCodec) Decode(c gnet.Conn) ([]byte, error) {
	buf, err := c.Peek(payloadOffset)
	if err != nil {
		if errors.Is(err, io.ErrShortBuffer) {
			return nil, ErrIncompletePacket
		}
		return nil, err
	}
	if !bytes.Equal(magicBytes, buf[:magicBytesSize]) {
		return nil, ErrMagicMismatch
	}
	payloadSize := binary.BigEndian.Uint32(buf[magicBytesSize:payloadOffset])
	totalSize := payloadOffset + int(payloadSize)
	buf, err = c.Peek(totalSize)
	if err != nil {
		if errors.Is(err, io.ErrShortBuffer) {
			return nil, ErrIncompletePacket
		}
		return nil, err
	}
	payload := make([]byte, payloadSize)
	copy(payload, buf[payloadOffset:totalSize])
	_, _ = c.Discard(totalSize)
	return payload, nil
}
