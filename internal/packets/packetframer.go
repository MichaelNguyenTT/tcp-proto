package packets

import (
	"fmt"
	"io"
)

type PacketFramer struct {
	buff       []byte
	pos        int
	outputChan chan *Packet
}

func NewPacketFramer() PacketFramer {
	return PacketFramer{
		buff:       make([]byte, MAX_PACKET_SIZE),
		outputChan: make(chan *Packet, 5),
	}
}

// TODO: finish framing packets
func FramePacketFromReader(framer *PacketFramer, r io.Reader) error {
	// tempbuf will consume 100 bytes at a time when during reading
	tempBuff := make([]byte, 100)
	for {
		n, err := r.Read(tempBuff)
		if err != nil {
			if err != io.EOF {
				break
			}

			return err
		}
		fmt.Printf("Read %d bytes from data: %v\n", n, tempBuff[:n])
		framer.InsertBuffer(tempBuff[:n])
	}

	return nil
}

func (p *PacketFramer) InsertBuffer(data []byte) error {
	n := copyPacketData(p.buff[:p.pos], data)

	// add the remaining data to the buffer
	if n < len(data) {
		p.buff = append(p.buff, data[n:]...)
	}

	for {
		cpkt, err := p.extractPackets()
		if err != nil {
			return err
		}
		// send complete packets to channel
		p.outputChan <- cpkt
	}
}

func (p *PacketFramer) extractPackets() (*Packet, error) {
	pktLen := extractPacketLength(p.buff)
	if pktLen > MAX_PACKET_SIZE {
		return nil, ErrPktExceededMax
	}

	// entire packet length
	fullpktLen := pktLen + HEADER_SIZE

	//TODO: make it better...feels too hacky lol
	if fullpktLen <= uint16(p.pos) {
		// out packet will be returned and sent to the channel
		outpkt := make([]byte, fullpktLen) // empty buffer size of full pkt length

		// copy the packet buffer from the main frame buffer
		copyPacketData(outpkt, p.buff[:fullpktLen])
		// remove the full packet data from the main frame buffer
		// basically moves the next packet to the front of the buffer
		removePacketData(p.buff, p.buff[fullpktLen:])

		// update the position in the frame
		p.pos = p.pos - int(fullpktLen)

		pkt := createPacketFromRawBytes(outpkt)

		return &pkt, nil
	}

	return nil, nil
}
