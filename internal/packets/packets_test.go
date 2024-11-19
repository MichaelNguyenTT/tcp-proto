package packets_test

import (
	"bytes"
	"encoding/binary"
	"tcpserver/internal/packets"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	HEADER_SIZE     = 3
	MAX_PACKET_SIZE = 1024
)

type TestPacket struct {
	pkt []byte
	len int
}

type PacketType uint8

const (
	PacketMessage PacketType = iota
)

var testData = []byte("hello world")

func (p *TestPacket) Read(data []byte) (int, error) {
	return copy(data, testData), nil
}

func (p *TestPacket) Encoder() ([]byte, error) {
	return p.pkt, nil
}

func (p *TestPacket) PacketType() uint8 {
	return uint8(PacketMessage)
}

func TestNewPacket(t *testing.T) {
	testpkt := &TestPacket{}

	pkt := packets.NewPacket(testpkt)

	buff := make([]byte, 0, 50)
	b := bytes.NewBuffer(buff)

	n, err := pkt.WritePacket(b)

	require.NoError(t, err, "packet logic had an error")
	require.Equal(t, n, HEADER_SIZE+len(testData), "n packet length should equal to testData")
	require.Equal(t, testData, buff[HEADER_SIZE:n], "data bytes should be equal to testData")
}

type testMessage struct {
	content []byte
	pType   PacketType
}

func (m *testMessage) Read(data []byte) (int, error) {
	return copy(data, m.content), nil
}

func (m *testMessage) PacketType() uint8 {
	return uint8(m.pType)
}

func (m *testMessage) Encoder() ([]byte, error) {
	return m.content, nil
}

func TestPacketEncoderFormat(t *testing.T) {
	payload := []byte{0x01, 0x02, 0x03}
	msg := &testMessage{
		content: payload,
		pType:   PacketMessage,
	}

	p := packets.NewPacket(msg)
	encoded, err := p.Encoder()

	expected := []byte{
		0x00, 0x03, // length (3)
		0x00,             // type (PacketMessage)
		0x01, 0x02, 0x03, // payload
	}

	require.NoError(t, err, "Encoder() error: %v")
	require.Equal(t, expected, encoded, "should have same encoding format as expected")
}

func TestPacketFromReader(t *testing.T) {
	data := []byte{
		0x00, 0x06,
		0x00,
		'P', 'A', 'C', 'K', 'E', 'T',
	}

	b := []byte("PACKET")

	pktSize := binary.BigEndian.Uint16(data)

	reader := bytes.NewReader(data)

	pkt, err := packets.PacketFromReader(reader)

	require.NoError(t, err, "packet reader error")
	require.Equal(t, int(pktSize), int(pkt.Len()))
	require.Equal(t, b, pkt.DataBytes(), "should have same data bytes")
}
