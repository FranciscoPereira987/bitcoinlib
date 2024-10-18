package bitcoinlib

import (
	"encoding/hex"
	"fmt"
	"net"
)

type SimpleNode struct {
	host       [16]byte
	port       uint16
	connection net.Conn
	testnet    bool
	logging    bool
}

type NodeParams struct {
	Addr    string
	Port    uint16
	Testnet bool
	Logging bool
}

func NewSimpleNode(params NodeParams) *SimpleNode {
	if params.Port == 0 && params.Testnet {
		params.Port = 18333
	} else if params.Port == 0 {
		params.Port = 8333
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", params.Addr, params.Port))
	if err != nil {
		panic(fmt.Sprintf("Could not connect to host %s because of %s", params.Addr, err))
	}
	return &SimpleNode{
		IPAddressFromString(conn.RemoteAddr().String()),
		params.Port,
		conn,
		params.Testnet,
		params.Logging,
	}
}

func (sn *SimpleNode) Send(message Message) error {
	envelope := NewNetworkMessage(sn.testnet)
	envelope.command = message.Command()
	envelope.payload = message.Serialize()
	serialized := envelope.Serialize()
	if sn.logging {
		fmt.Printf("Sending: %s\n", hex.EncodeToString(serialized))
	}
	_, err := sn.connection.Write(serialized)
	return err
}

func (sn *SimpleNode) Read() (*NetworkMessage, error) {
	message := NewNetworkMessage(sn.testnet)
	err := message.Parse(sn.connection)
	if sn.logging {
		fmt.Printf("Recieved: %o\nWith Error: %s\n", message, err)
	}
	return message, err
}

func (sn *SimpleNode) WaitFor(commands map[string]Message) (Message, error) {
	var command string
	var envelope *NetworkMessage
	var err error
	for _, ok := commands[command]; !ok; _, ok = commands[command] {
		envelope, err = sn.Read()
		if err != nil {
			return nil, err
		}
		fmt.Println("Read a message")
		command = envelope.GetCommand()
		if envelope.EqCommand(VERSION_MESSAGE) {
			sn.Send(NewVerackMessage())
		} else if envelope.EqCommand(PING_MESSAGE) {
			PING_MESSAGE.Parse(envelope.payload)
			PONG_MESSAGE.nonce = PING_MESSAGE.nonce
			sn.Send(PONG_MESSAGE)
		}
	}
	return commands[command].Parse(envelope.payload)
}
