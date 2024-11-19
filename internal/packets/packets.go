package packets

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	HEADER_SIZE            = 3
	MAX_PACKET_SIZE        = 1024
	PACKET_LENGTH_POSITION = 0
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

	ErrIncompleteHeader = fmt.Errorf("incomplete header read")
	ErrIncompleteData   = fmt.Errorf("incomplete data read")
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

func (p *Packet) Len() uint16 {
	return binary.BigEndian.Uint16(p.pkt[PACKET_LENGTH_POSITION:])
}

func (p *Packet) DataBytes() []byte {
	return p.pkt[HEADER_SIZE:]
}

func (p *Packet) ToString() string {
	return fmt.Sprintf("Packet length: %d\n Data: %s\n", p.len, p.pkt[HEADER_SIZE:])
}

func (p *Packet) Encoder() ([]byte, error) {
	pkt := p.createPacketBuffer()

	if len(p.pkt) < len(pkt) {
		return nil, ErrPktExceededMax
	}

	dataLeng := p.len - HEADER_SIZE

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

func PacketFromReader(reader io.Reader) (*Packet, error) {
	header := make([]byte, HEADER_SIZE)
	n, err := reader.Read(header)
	if err != nil {
		return nil, err
	}

	if n < HEADER_SIZE {
		return nil, ErrIncompleteHeader
	}

	pktSize := extractPacketLength(header)

	fullPkt := int(pktSize) + HEADER_SIZE
	if fullPkt > MAX_PACKET_SIZE {
		return nil, ErrPktExceededMax
	}

	// create buffer for complete packet
	pktBuf := make([]byte, fullPkt, fullPkt)

	copy(pktBuf, header) // copy header

	data := pktBuf[HEADER_SIZE:]
	n, err = reader.Read(data)
	if err != nil {
		return nil, err
	}

	if n < int(pktSize) {
		return nil, ErrIncompleteData
	}

	pkt := createPacketFromRawBytes(pktBuf)
	return &pkt, nil
}
