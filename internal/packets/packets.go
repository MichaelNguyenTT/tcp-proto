package packets

import "io"

const (
	HEADER_SIZE     = 3
	MAX_PACKET_SIZE = 1024
)

type PacketType uint8

const (
	PacketMessage PacketType = iota
	PacketError
)

type Encoding uint8

const (
	EncodeBytes Encoding = iota
	EncodeString
)

type Packet struct {
	data []byte
	len  int
}

// encoder to write the packet
type PacketEncoder interface {
	io.Reader
	Type() uint8
	Encoder() Encoding
}

// packet protocol structure
// [type: 1 byte][packet size: 2 bytes][data: number of bytes]

// packet type = uint8 1 byte for memory efficiency
// packet size = length of the data stored with byte order as uint16, taking 2 bytes

//TODO: packet logic
