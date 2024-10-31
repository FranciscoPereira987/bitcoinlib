package bitcoinlib

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"
)

const MAINNET_MAGIC = 0xf9beb4d9
const TESTNET_MAGIC = 0x0b110907

const VERACK = "verack"
const VERSION = "version"
const PING = "ping"
const PONG = "pong"
const HEADERS = "headers"
const GETHEADERS = "getheaders"
const MERKLEBLOCK = "merkleblock"

var IPV4_BASE = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0, 0, 0, 0}

var VERACK_COMMAND = IntoCommand(VERACK)
var VERSION_COMMAND = IntoCommand(VERSION)
var PING_COMMAND = IntoCommand(PING)
var PONG_COMMAND = IntoCommand(PONG)
var GETHEADERS_COMMAND = IntoCommand(GETHEADERS)
var HEADERS_COMMAND = IntoCommand(HEADERS)
var MERKLEBLOCK_COMMAND = IntoCommand(MERKLEBLOCK)

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

type GetHeadersMessage struct {
	version    uint32
	hashes     uint8
	startBlock string
	endBlock   string
}

type HeadersMessage struct {
	blocks []*Block
}

type MerkleBlockMessage struct {
	blocks            *Block
	flagBits          []byte
	totalTransactions uint32
}

func NewGetHeadersMessage(startBlock string, endBlock string) *GetHeadersMessage {
	return &GetHeadersMessage{
		version:    70015,
		hashes:     1,
		startBlock: startBlock,
		endBlock:   endBlock,
	}
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

func NewHeadersMessage() *HeadersMessage {
	return &HeadersMessage{make([]*Block, 0)}
}

func NewMerkleBlockMessage() *MerkleBlockMessage {
	return &MerkleBlockMessage{}
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
	actual := total
	for total != 0 && err == nil && actual < len(buf)-1 {
		total, err = from.Read(buf[actual:])
		actual += total
	}
	if err != nil || actual < len(buf)-1 {
		err = errors.Join(err, fmt.Errorf("payload invalid read: %d vs %d", total, len(buf)))
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

func (m *GetHeadersMessage) Serialize() []byte {
	buf := binary.LittleEndian.AppendUint32(nil, m.version)
	buf = append(buf, m.hashes)
	start, err := hex.DecodeString(m.startBlock)
	if err != nil || len(start) != 32 {
		panic(fmt.Sprintf("Error: %s with length %d", err, len(start)))
	}
	end, err := hex.DecodeString(m.endBlock)
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
	if len(end) != 32 {
		end = make([]byte, 32)
	}
	//Need to use start and end as little endian
	slices.Reverse(start)
	slices.Reverse(end)
	buf = append(buf, start...)
	return append(buf, end...)
}

func (m *GetHeadersMessage) Command() [12]byte {
	return GETHEADERS_COMMAND
}

func (m *GetHeadersMessage) Parse(stream []byte) (Message, error) {
	m.version = binary.LittleEndian.Uint32(stream)
	stream = stream[4:]
	m.hashes = stream[0]
	stream = stream[1:]
	if len(stream) != 64 {
		return nil, fmt.Errorf("invalid stream length for blocks id, expected 64, but got %d", len(stream))
	}
	start := stream[:32]
	slices.Reverse(start)
	end := stream[32:]
	slices.Reverse(end)
	m.startBlock = hex.EncodeToString(start)
	m.endBlock = hex.EncodeToString(end)
	return m, nil
}

func (m *HeadersMessage) Serialize() []byte {
	buf := binary.AppendUvarint(nil, uint64(len(m.blocks)))
	for _, block := range m.blocks {
		buf = append(buf, block.Serialize()...)
		buf = append(buf, 0) //Append number of transactions
	}
	return buf
}

func (m *HeadersMessage) Command() [12]byte {
	return HEADERS_COMMAND
}

func (m *HeadersMessage) Parse(stream []byte) (Message, error) {
	_, read := binary.Varint(stream)
	total := ReadVarInt(bytes.NewReader(stream[:read]))
	stream = stream[read:]
	m.blocks = make([]*Block, total)
	readStream := bytes.NewReader(stream)
	for i := range total {
		block := NewBlock()
		err := block.Parse(readStream)
		if err != nil {
			return nil, err
		}
		read, err := readStream.ReadByte()
		if err != nil {
			return nil, err
		}
		if read != 0 {
			return nil, fmt.Errorf("Read total transactions != 0 (%d)", read)
		}
		m.blocks[i] = block
	}
	return m, nil
}

func (m *HeadersMessage) TotalBlocks() int {
	return len(m.blocks)
}

func (m *HeadersMessage) GetBlock(i int) *Block {
	return m.blocks[i]
}

func (m *MerkleBlockMessage) Command() [12]byte {
	return MERKLEBLOCK_COMMAND
}

func (m *MerkleBlockMessage) Serialize() []byte {
	if m.blocks == nil {
		return []byte{}
	}
	buf := m.blocks.Serialize()
	buf = binary.LittleEndian.AppendUint32(buf, m.totalTransactions)
	buf = append(buf, m.blocks.SerializeHashes()...)
	totalBits := len(m.flagBits)
	buf = append(buf, binary.LittleEndian.AppendUint16(nil, uint16(totalBits))[0])
	buf = append(buf, m.flagBits...)
	return buf
}

func (m *MerkleBlockMessage) Parse(stream []byte) (Message, error) {
	if m.blocks == nil {
		m.blocks = NewBlock()
	}
	buf := bytes.NewBuffer(stream)
	if err := m.blocks.Parse(buf); err != nil {
		return nil, err
	}
	transactions := make([]byte, 4)
	if _, err := buf.Read(transactions); err != nil {
		return nil, err
	}
	m.totalTransactions = binary.LittleEndian.Uint32(transactions)
	hashes := ReadVarInt(buf)
	hashesBuf := make([][]byte, hashes)
	for i := range hashes {
		helper := make([]byte, 32)
		if _, err := buf.Read(helper); err != nil {
			fmt.Printf("Going out through here: %d from %d\n", i, hashes)
			return nil, err
		}
		hashesBuf[i] = helper
	}
	m.blocks.hashes = hashesBuf
	totalBytes := []byte{0}
	if _, err := buf.Read(totalBytes); err != nil {
		return nil, err
	}
	flagBytes := make([]byte, totalBytes[0])
	if _, err := buf.Read(flagBytes); err != nil {
		return nil, err
	}
	m.flagBits = flagBytes
	return m, nil
}

func (m *MerkleBlockMessage) FlagsToBits() []bool {
	flags := make([]bool, 0)
	for _, byte := range m.flagBits {
		for _ = range 8 {
			flags = append(flags, byte&1 == 1)
			byte = byte >> 1
		}
	}
	return flags
}

func (m *MerkleBlockMessage) ValidateTree() bool {
	flags := m.FlagsToBits()
	return m.blocks.ValidateMerkle(flags, int(m.totalTransactions))
}
