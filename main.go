package main

import (
	"bitcoinlib/bitcoinlib"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"bytes"
	"encoding/hex"
	"fmt"
)

func main() {
	passkey := "41797243 Francisco Javier Pereira Argentina!"
	hashed := bitcoinlib.Hash256([]byte(passkey))
	key := bitcoinlib.NewPrivateKey(bitcoinlib.FromHexString(hex.EncodeToString(hashed)))
	txString, _ := hex.DecodeString("01000000022f2afe57bde0822c793604baae834f2cd26155bf1c0d37480212c107e75cd011010000006a47304402204cc5fe11b2b025f8fc9f6073b5e3942883bbba266b71751068badeb8f11f0364022070178363f5dea4149581a4b9b9dbad91ec1fd990e3fa14f9de3ccb421fa5b269012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67ffffffff153db0202de27e7944c7fd651ec1d0fab1f1aaed4b0da60d9a1b06bd771ff651010000006b483045022100b7a938d4679aa7271f0d32d83b61a85eb0180cf1261d44feaad23dfd9799dafb02205ff2f366ddd9555f7146861a8298b7636be8b292090a224c5dc84268480d8be1012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67ffffffff01d0754100000000001976a914ad346f8eb57dee9a37981716e498120ae80e44f788ac00000000")
	otherTx, _ := bitcoinlib.ParseTransaction(bytes.NewReader(txString))
	fmt.Println(otherTx.Verify(true))
	//Create the transaction
	tx := bitcoinlib.NewTransaction()
	
	//Add the input
	tx.AddInput("7157c98038b093050feef0899a6b6bad7f99c3129b4c56fb44bd0efcd9bbf542", 1)
	tx.AddInput("8071fbf1198849bd643259fa86d758e3b78ce649eb8ffb0cee3107344f8566aa",0)
	//Add the outputs
	tx.AddOutput(20000,"mwQkTVnb1hLa6qXyLT3i2cAFmi8p8Wn5wr")
	//Sign the transaction
	tx.Sign(true, key)
	fmt.Println(tx.Fee(true))
	rawTx, err := bitcoin.TxFromHex(hex.EncodeToString(tx.Serialize()))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	//If I tried to broadcast it using the hex method, it wouldnt work
	//But with the rawTx it would even though they are equal
	fmt.Println(rawTx.String() == hex.EncodeToString(tx.Serialize())) 
	//Verify it
	fmt.Println(tx.Verify(true))
	//Print its serialization
	//ae55a2b58fd4839a5e597d94fa4c80c6195ada82
	fmt.Println(key.Address(bitcoinlib.COMPRESSED, true))
}
