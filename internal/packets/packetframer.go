package packets

import (
	"fmt"
	"io"
)

type PacketFramer struct {
	buff    []byte
	pos     int
	PacketC chan *Packet
}

func NewPacketFramer() PacketFramer {
	return PacketFramer{
		buff:    make([]byte, MAX_PACKET_SIZE),
		PacketC: make(chan *Packet, 5),
	}
}

// TODO: finish framing packets
func FramePacketFromReader(pkt *PacketFramer, r io.Reader) error {
	// tempbuf will consume 100 bytes at a time when during reading
	tempBuf := make([]byte, 100)
	for {
		// read the number of bytes from stream of data
		n, err := r.Read(tempBuf)
		if err != nil {
			if err != io.EOF {
				break
			}

			return err
		}

		fmt.Printf("Read %d bytes from data: %v\n", n, tempBuf[:n])
	}

	return nil
}
