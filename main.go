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
	tx.AddInput("ee3f743e3cba5ddb75cdf77cfdfaddaeb2ce00ad8c7a92b9338cf4bc05c7db28", 0)
	tx.AddInput("2fa03ac8be24b9ee984737130694aeed71d0a737d36c896a3d2ed898461aaa25",0)
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
