package packets

import "encoding/binary"

// returns the first 2 bytes from the header
func extractPacketLength(data []byte) uint16 {
	pktLength := binary.BigEndian.Uint16(data[0:2])
	return pktLength
}

func createPacketFromRawBytes(data []byte) Packet {
	// need some sort of assertions for simulation testing
	// length := len(data) - HEADER_SIZE
	// pktLen := extractPacketLength(data)

	return Packet{
		pkt: data,
		len: len(data),
	}
}

func copyPacketData[T any](to, from []T) int {
	return copy(to, from)
}

func removePacketData[T any](to, from []T) int {
	return copy(to, from)
}
