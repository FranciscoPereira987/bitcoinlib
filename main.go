package main

import (
	"bitcoinlib/bitcoinlib"
	"encoding/hex"
	"fmt"
)

func main() {
	passkey := "41797243 Francisco Javier Pereira Argentina!"
	hashed := bitcoinlib.Hash256([]byte(passkey))
	key := bitcoinlib.NewPrivateKey(bitcoinlib.FromHexString(hex.EncodeToString(hashed)))
	//Create the transaction
	tx := bitcoinlib.NewTransaction()
	//Add the input
	tx.AddInput("e89f524d861e8165d41f6fd6a482041f976214c2a3be1b00fc7e3e72473a0991", 0)
	//Add the outputs
	tx.AddOutput(8952, "mwJn1YPMq7y5F8J3LkC5Hxg9PHyZ5K4cFv")
	tx.AddOutput(5800, "mwQkTVnb1hLa6qXyLT3i2cAFmi8p8Wn5wr")
	//Sign the transaction
	tx.SignInput(0, true, key)
	//Verify it
	fmt.Println(tx.Verify(true))
	//Print its serialization
	//ae55a2b58fd4839a5e597d94fa4c80c6195ada82
	fmt.Println(hex.EncodeToString(bitcoinlib.Hash160(key.Sec(bitcoinlib.COMPRESSED))))
	fmt.Println(key.Address(bitcoinlib.COMPRESSED, true))
}
