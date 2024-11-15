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
	pkt []byte
	len int
}

type PacketEncoder interface {
	io.Reader
	PacketType() uint8
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

// new packet tester
func NewPacket(pkt PacketEncoder) *Packet {
	buf := make([]byte, MAX_PACKET_SIZE)
	data := buf[HEADER_SIZE:]
	n, _ := pkt.Read(data)

	return &Packet{
		pkt: buf,
		len: HEADER_SIZE + n,
	}
}

func (p *Packet) PacketType() PacketType {
	return PacketType(p.pkt[2])
}

func (p *Packet) Read(data []byte) (int, error) {
	if len(data) < p.len {
		return 0, ErrPktTooSmall
	}

	copy(data, p.pkt[:p.len])

	return p.len, nil
}

func (p *Packet) createPacketBuffer() []byte {
	buf := make([]byte, p.len)
	return buf
}

func (p *Packet) getPacketDataLength() int {
	return p.len - HEADER_SIZE
}

func (p *Packet) Encoder() ([]byte, error) {
	pkt := p.createPacketBuffer()

	if len(p.pkt) < len(pkt) {
		return nil, ErrPktExceededMax
	}

	dataLeng := p.getPacketDataLength()

	binary.BigEndian.PutUint16(pkt[0:2], uint16(dataLeng))
	pkt[2] = byte(p.PacketType())

	copy(pkt[HEADER_SIZE:], p.pkt[HEADER_SIZE:p.len])

	return pkt, nil
}

func (p *Packet) WritePacket(w io.Writer) (int, error) {
	if p.len > MAX_PACKET_SIZE {
		return 0, ErrPktExceededMax
	}

	n, err := w.Write(p.pkt[:p.len])
	if err != nil {
		return 0, nil
	}

	return n, nil
}
