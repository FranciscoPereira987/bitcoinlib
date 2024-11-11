package main

import (
	"bitcoinlib/bitcoinlib"
	"bytes"
	"encoding/hex"
	"fmt"
)

func transactionMain() {
	passkey := "41797243 Francisco Javier Pereira Argentina!"
	hashed := bitcoinlib.Hash256([]byte(passkey))
	key := bitcoinlib.NewPrivateKey(bitcoinlib.FromHexString(hex.EncodeToString(hashed)))
	//Create the transaction
	tx := bitcoinlib.NewTransaction()
	//Add the input
	tx.AddInput("ee3f743e3cba5ddb75cdf77cfdfaddaeb2ce00ad8c7a92b9338cf4bc05c7db28", 0)
	tx.AddInput("2fa03ac8be24b9ee984737130694aeed71d0a737d36c896a3d2ed898461aaa25", 0)
	//Add the outputs
	tx.AddOutput(30000, "mwQkTVnb1hLa6qXyLT3i2cAFmi8p8Wn5wr")
	//Sign the transaction
	tx.Sign(true, key)
	fmt.Println(tx.Fee(true))
	fmt.Println(tx.String())
	//Verify it
	fmt.Println(tx.Verify(true))
	//Print its serialization
	//ae55a2b58fd4839a5e597d94fa4c80c6195ada82
	fmt.Println(key.Address(bitcoinlib.COMPRESSED, true))
}

func nodeMain() {
	params := bitcoinlib.NodeParams{
		Addr:    "testnet-seed.bitcoin.jonasschnelli.ch",
		Testnet: true,
	}
	node := bitcoinlib.NewSimpleNode(params)
	err := node.Handshake()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func nodeHeaders() {
	params := bitcoinlib.NodeParams{
		Addr:    "testnet-seed.bitcoin.jonasschnelli.ch",
		Testnet: true,
	}
	gBytes, _ := hex.DecodeString(bitcoinlib.TESTNET_GENESIS_BLOCK)
	block := bitcoinlib.NewBlock()
	block.Parse(bytes.NewReader(gBytes))
	node := bitcoinlib.NewSimpleNode(params)
	err := node.Handshake()
	if err != nil {
		fmt.Printf("Error during handshake: %s\n", err)
		return
	}
	node.Send(bitcoinlib.NewGetHeadersMessage(block.Hash(), ""))
	maped := map[string]bitcoinlib.Message{
		bitcoinlib.HEADERS: bitcoinlib.NewHeadersMessage(),
	}
	result, err := node.WaitFor(maped)
	if err != nil {
		fmt.Printf("Error waiting for headers message: %s\n", err)
		return
	}
	headers := result.(*bitcoinlib.HeadersMessage)
	for i := 0; i < headers.TotalBlocks() && headers.TotalBlocks() > 0; {
		block = headers.GetBlock(i)
		if !block.CheckPOW() {
			fmt.Printf("Failed checking block %d POW\n", i)
		}
		if i%2016 == 0 {
			fmt.Printf("Difficulty: %s\n", block.Difficulty())
		}
		if i == headers.TotalBlocks()-1 {
			fmt.Println("Asking for more")
			node.Send(bitcoinlib.NewGetHeadersMessage(block.Hash(), ""))
			result, err := node.WaitFor(maped)
			if err != nil {
				fmt.Printf("Error waiting for headers message: %s\n", err)
				return
			}
			headers = result.(*bitcoinlib.HeadersMessage)
			fmt.Println("Total blocks: ", headers.TotalBlocks())
			i = -1
		}
		i++
	}

}

func TransactionOfInterestMain() {
	startBlock := "00000000000538d5c2246336644f9a4956551afb44ba47278759ec55ea912e19"
	address := "mwJn1YPMq7y5F8J3LkC5Hxg9PHyZ5K4cFv"
	params := bitcoinlib.NodeParams{
		Addr:    "testnet-seed.bitcoin.jonasschnelli.ch",
		Testnet: true,
		Logging: true,
	}
	node := bitcoinlib.NewSimpleNode(params)
	if err := node.Handshake(); err != nil {
		fmt.Printf("Failed handshake: %s", err)
		return
	}
	h160 := bitcoinlib.FromBase58Address(address)
	h160Hex, _ := hex.DecodeString(h160)
	filter := bitcoinlib.NewBloomFilter(30)
	filterParams := &bitcoinlib.MurmurParams{
		FunctionCount: 5,
		Tweak:         90210,
	}
	filter.Set(h160Hex, filterParams)
	node.Send(&bitcoinlib.FilterLoadMessage{Filter: filter})
	headers := bitcoinlib.NewGetHeadersMessage(startBlock, "")
	node.Send(headers)
	_, err := node.WaitFor(map[string]bitcoinlib.Message{bitcoinlib.HEADERS: bitcoinlib.NewHeadersMessage()})
	if err != nil {
		fmt.Printf("Failed recovering headers message: %s", err)
		return
	}
	node.WaitFor(map[string]bitcoinlib.Message{bitcoinlib.PING: bitcoinlib.PING_MESSAGE})
}

func main() {
	TransactionOfInterestMain()
}
