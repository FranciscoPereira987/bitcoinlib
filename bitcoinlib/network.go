package bitcoinlib

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

const MAINNET_MAGIC = 0xf9beb4d9
const TESTNET_MAGIC = 0x0b110907

const VERACK = "verack"
const VERSION = "version"
const PING = "ping"
const PONG = "pong"

var IPV4_BASE = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0, 0, 0, 0}

var VERACK_COMMAND = IntoCommand(VERACK)
var VERSION_COMMAND = IntoCommand(VERSION)
var PING_COMMAND = IntoCommand(PING)
var PONG_COMMAND = IntoCommand(PONG)

var VERSION_MESSAGE = NewVersionMessage()
var VERACK_MESSAGE = NewVerackMessage()
var PONG_MESSAGE = NewPongMessage(0)
var PING_MESSAGE = NewPingMessage(0)

func IPAddressFromString(add string) [16]byte {
	converted, _ := hex.DecodeString(add)
	var ipAddr [16]byte = IPV4_BASE
	if len(converted) == 4 {
		copy(ipAddr[12:], converted)
	} else {
		//I assume its an IPV6 Addr
		copy(ipAddr[:], converted)
	}
	return ipAddr
}

func IntoCommand(value string) [12]byte {
	buf := [12]byte{}
	if len(value) <= 12 {
		copy(buf[:], []byte(value))
	}
	return buf
}

type Message interface {
	Serialize() []byte
	Command() [12]byte
	Parse([]byte) (Message, error)
}

type NetworkMessage struct {
	magic   uint32
	command [12]byte
	payload []byte
}

type VersionMessage struct {
	Protocol         uint32
	Services         uint64
	Timestamp        uint64
	RecieverServices uint64
	RecieverAddress  [16]byte
	RecieverPort     uint16
	SenderServices   uint64
	SenderAddress    [16]byte
	SenderPort       uint16
	Nonce            uint64
	UserAgent        string
	Height           uint32
	RelayFlag        bool
}

type VerackMessage struct{}

type PingMessage struct {
	nonce uint64
}

type PongMessage struct {
	nonce uint64
}

func NewPingMessage(nonce uint64) *PingMessage {
	return &PingMessage{nonce}
}

func NewPongMessage(nonce uint64) *PongMessage {
	return &PongMessage{nonce}
}

func NewVerackMessage() *VerackMessage {
	return &VerackMessage{}
}

func NewVersionMessage() *VersionMessage {
	time := uint64(time.Now().Unix())
	nonce, _ := rand.Int(rand.Reader, FromInt(2).Exp(FromInt(64), ZERO).value)
	return &VersionMessage{
		Protocol:         70015,
		Services:         0,
		Timestamp:        time,
		RecieverServices: 0,
		RecieverAddress:  IPV4_BASE,
		RecieverPort:     8333,
		SenderServices:   0,
		SenderAddress:    IPV4_BASE,
		SenderPort:       8333,
		Nonce:            uint64(nonce.Int64()),
		UserAgent:        "/programmingbitcoin:0.1/",
		RelayFlag:        false,
	}
}

func NewNetworkMessage(testnet bool) *NetworkMessage {
	var magic uint32
	if testnet {
		magic = TESTNET_MAGIC 
	} else {
		magic = MAINNET_MAGIC
	}
	return &NetworkMessage{
		magic,
		[12]byte{},
		nil,
	}
}

func (nm *NetworkMessage) EqCommand(m Message) bool {
	other := m.Command()
	return string(nm.command[:]) == string(other[:])
}

func readMagic(from io.Reader) (magic uint32, err error) {
	buf := make([]byte, 4)
	total, err := from.Read(buf)
	if err != nil || total < len(buf) {
    err = errors.Join(err, errors.New("invalid magic read"))
	} else {
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
	buf = make([]byte, 4+binary.LittleEndian.Uint32(buf))
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
	for lastIndex < len(m.command) && m.command[lastIndex] != 0 {
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
	buf = binary.LittleEndian.AppendUint64(buf, m.SenderServices)
	buf = append(buf, m.SenderAddress[:]...)
	buf = binary.BigEndian.AppendUint16(buf, m.SenderPort)
	buf = binary.BigEndian.AppendUint64(buf, m.Nonce)
	buf = append(buf, EncodeVarInt(uint64(len(m.UserAgent)))...)
	buf = append(buf, []byte(m.UserAgent)...)
	buf = binary.LittleEndian.AppendUint32(buf, m.Height)
	if m.RelayFlag {
		buf = append(buf, 1)
	} else {
		buf = append(buf, 0)
	}
	return buf
}

func (m *VersionMessage) Command() [12]byte {
	return VERSION_COMMAND
}

func (m *VerackMessage) Serialize() []byte {
	return []byte{}
}

func (m *VerackMessage) Command() [12]byte {
	return VERACK_COMMAND
}

func (m *VersionMessage) Parse(stream []byte) (Message, error) {
  m.Protocol = binary.LittleEndian.Uint32(stream)
	stream = stream[4:]
	m.Services = binary.LittleEndian.Uint64(stream)
	stream = stream[8:]
	m.Timestamp = binary.LittleEndian.Uint64(stream)
	stream = stream[8:]
	m.RecieverServices = binary.LittleEndian.Uint64(stream)
	stream = stream[8:]
	m.RecieverAddress = [16]byte(stream[:16])
	stream = stream[16:]
	m.RecieverPort = binary.BigEndian.Uint16(stream)
	stream = stream[2:]
	m.SenderServices = binary.LittleEndian.Uint64(stream)
	stream = stream[8:]
	m.SenderAddress = [16]byte(stream[:16])
	stream = stream[16:]
	m.SenderPort = binary.BigEndian.Uint16(stream)
	stream = stream[2:]
	m.Nonce = binary.LittleEndian.Uint64(stream)
	stream = stream[8:]
	userAgentLength, total := binary.Varint(stream) 
	stream = stream[total:]
	m.UserAgent = string(stream[:userAgentLength])
	stream = stream[userAgentLength:]
	m.Height = binary.BigEndian.Uint32(stream)
	stream = stream[4:]
	m.RelayFlag = (len(stream) > 0 && stream[0] == 1)
	return m, nil
}

func (m *VerackMessage) Parse(stream []byte) (Message, error) {
	if len(stream) != 0 {
		return nil, errors.New("invalid Verack payload")
	}
	return m, nil
}

func (m *PongMessage) Serialize() []byte {
	return binary.BigEndian.AppendUint64(nil, m.nonce)
}

func (m *PongMessage) Command() [12]byte {
	return PONG_COMMAND
}

func (m *PongMessage) Parse(stream []byte) (Message, error) {
	if len(stream) != 8 {
		return nil, errors.New("invalid poing message")
	}
	m.nonce = binary.BigEndian.Uint64(stream)
	return m, nil
}

func (m *PingMessage) Serialize() []byte {
	return binary.BigEndian.AppendUint64(nil, m.nonce)
}

func (m *PingMessage) Command() [12]byte {
	return PING_COMMAND
}

func (m *PingMessage) Parse(stream []byte) (Message, error) {
	if len(stream) != 8 {
		return nil, errors.New("invalid poing message")
	}
	m.nonce = binary.BigEndian.Uint64(stream)
	return m, nil
}
