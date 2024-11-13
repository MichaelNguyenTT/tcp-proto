package packets

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	HEADER_SIZE     = 3
	MAX_PACKET_SIZE = 1024
)

type PacketType uint8

const (
	PacketMessage PacketType = iota
	PacketError
)

type Packet struct {
	data   []byte
	pktLen int
}

type PacketEncoder interface {
	io.Reader
	Type() uint8
	Encoder() ([]byte, error)
}

var (
	ErrPktExceededMax = fmt.Errorf("buffer exceeded maximum packet size: %d", MAX_PACKET_SIZE)
	ErrPktTooSmall    = fmt.Errorf("buffer data too small for packet")

	ErrInvalidHeader = fmt.Errorf("invalid packet header")
	ErrInvalidType   = fmt.Errorf("invalid packet type")
)

// packet protocol structure
// [packet size: 2 bytes][type: 1 bytes][data: n of bytes]

//TODO: packet logic

func NewPacket(pkt PacketEncoder) Packet {
	buf := make([]byte, MAX_PACKET_SIZE)

	data := buf[HEADER_SIZE:]
	n, _ := pkt.Read(data)

	binary.BigEndian.PutUint16(buf[0:2], uint16(n))
	buf[2] = byte(PacketType(PacketMessage))

	return Packet{
		data:   buf,             // data is the entire buffered packet
		pktLen: HEADER_SIZE + n, // pktLen => is basically length of byte slice of the payload + header size protocol || 4 + N
	}
}

// allocate empty buffer size with HEADER + data
func allocateBuffer(data []byte) []byte {
	pktSize := make([]byte, HEADER_SIZE+len(data))
	return pktSize
}

// HACK: feels wonky....
func (p *Packet) Type() PacketType {
	return PacketType(p.data[2])
}

// encoder will transform byte sequence to transmittable bytes
func (p *Packet) Encoder() ([]byte, error) {
	pktBuf := allocateBuffer(p.data)

	binary.BigEndian.PutUint16(pktBuf[0:2], uint16(p.pktLen))
	pktBuf[2] = byte(p.Type())

	copy(pktBuf[HEADER_SIZE:], p.data)

	return pktBuf, nil
}

func (p *Packet) Read(data []byte) (int, error) {
	if len(data) < p.pktLen {
		return 0, ErrPktTooSmall
	}

	copy(data, p.data[:p.pktLen])

	return p.pktLen, nil
}
