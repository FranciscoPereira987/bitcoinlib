package bitcoinlib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const MAINNET_MAGIC = 0xf9beb4d9
const TESTNET_MAGIC = 0X0b110907

var IPV4_BASE = [16]byte{0,0,0,0,0,0,0,0,0,0,0xff,0xff,0,0,0,0}

type NetworkMessage struct {
	magic uint32
	command [12]byte
	payload []byte
}

type VersionMessage struct {
	Protocol uint32
	Services uint64
	Timestamp uint64
	RecieverServices uint64
	RecieverAddress [16]byte
	RecieverPort uint16
	SenderServices uint64
	SenderAddress [16]byte
	SenderPort uint16
	Nonce uint64
	UserAgent string
	Height uint32
	RelayFlag bool
}

func NewVersionMessage() *VersionMessage {
	return &VersionMessage{
			Protocol: 70015,
			Services: 0,
			Timestamp: 0,
			RecieverServices: 0,
			RecieverAddress: IPV4_BASE,
			RecieverPort: 8333,
			SenderServices: 0,
			SenderAddress: IPV4_BASE,
			SenderPort: 8333,
			Nonce: 0,
			UserAgent: "/programmingbitcoin:0.1/",
			RelayFlag: false,
	}
}

func NewNetworkMessage(testnet bool) *NetworkMessage {
	var magic uint32 
	if testnet {
		magic = MAINNET_MAGIC
	}else {
		magic = TESTNET_MAGIC
	}
	return &NetworkMessage{
		magic,
		[12]byte{},
		nil,
	}
}

func readMagic(from io.Reader) (magic uint32, err error) {
	buf := make([]byte, 4)
	total, err := from.Read(buf)
	if err != nil || total < len(buf) {
		err = errors.Join(err, errors.New("invalid magic read"))
	}else {
		magic = binary.BigEndian.Uint32(buf)
	}
	return
}

func readCommand(from io.Reader) (command [12]byte, err error) {
	total, err := from.Read(command[:])
	if err != nil || total < len(command) {
		err = errors.Join(err, errors.New("invalid command read"))
	}
	return
}

func readPayload(from io.Reader) (checksum uint32, payload []byte, err error) {
	buf := make([]byte, 4)
	total, err := from.Read(buf)
	if err != nil || total < len(buf) {
		err = errors.Join(err, errors.New("invalid payload length read"))
		return
	}
	buf = make([]byte, 4 + binary.LittleEndian.Uint32(buf))
	total, err = from.Read(buf)
	if err != nil || total < len(buf) {
		err = errors.Join(err, errors.New("payload invalid read"))
	}
	checksum = binary.BigEndian.Uint32(buf)
	buf = buf[4:]
	payload = buf
	return
}

func checkChecksum(checksum uint32, payload []byte) (err error) {
	hashed := binary.BigEndian.Uint32(Hash256(payload))
	if hashed != checksum {
		err = fmt.Errorf("checksum did not match: %x vs %x", hashed, checksum)
	}
	return
}

func (m *NetworkMessage) Parse(from io.Reader) error {
	magic, err := readMagic(from)
	if err != nil {
		return err
	}
	command, err := readCommand(from)
	if err != nil {
		return err
	}
	checksum, payload, err := readPayload(from)
	if err != nil {
		return err
	}
	err = checkChecksum(checksum, payload)
	if err != nil {
		return err
	}
	m.magic = magic
	m.command = command
	m.payload = payload
	return nil
}

func (m *NetworkMessage) Serialize() []byte {
	buf := binary.BigEndian.AppendUint32(nil, m.magic)
	buf = append(buf, m.command[:]...)
	checksum := Hash256(m.payload)[:4]
	buf = binary.LittleEndian.AppendUint32(buf, uint32(len(m.payload)))
	buf = append(buf, checksum...)
	return append(buf, m.payload...)
}

func (m *NetworkMessage) GetCommand() string {
	lastIndex := 0
	for lastIndex < len(m.command) && m.command[lastIndex] != 0{
		lastIndex++
	}
	return string(m.command[:lastIndex])
}

func (m *NetworkMessage) GetPayload() []byte {
	return m.payload
}

func (m *VersionMessage) Serialize() []byte {
	buf := binary.LittleEndian.AppendUint32(nil, m.Protocol)
	buf = binary.LittleEndian.AppendUint64(buf, m.Services)
	buf = binary.LittleEndian.AppendUint64(buf, m.Timestamp)
	buf = binary.LittleEndian.AppendUint64(buf, m.RecieverServices)
	buf = append(buf, m.RecieverAddress[:]...)
	buf = binary.BigEndian.AppendUint16(buf, m.RecieverPort)
	buf = binary.LittleEndian.AppendUint64(buf,m.SenderServices)
	buf = append(buf, m.SenderAddress[:]...)
	buf = binary.BigEndian.AppendUint16(buf, m.SenderPort)
	buf = binary.BigEndian.AppendUint64(buf, m.Nonce)
	buf = append(buf, EncodeVarInt(uint64(len(m.UserAgent)))...)
	buf = append(buf, []byte(m.UserAgent)...)
	buf = binary.BigEndian.AppendUint32(buf, m.Height)
	if m.RelayFlag {
		buf = append(buf, 1)
	}else {
		buf = append(buf, 0)
	}
	return buf
}