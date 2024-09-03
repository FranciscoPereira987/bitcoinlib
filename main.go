package main

import (
	"bitcoinlib/bitcoinlib"
	"encoding/hex"
	"fmt"
)

func main() {
	passkey := "RandomPasskey"
	hashed := bitcoinlib.Hash256([]byte(passkey))
	key := bitcoinlib.NewPrivateKey(bitcoinlib.FromHexString(hex.EncodeToString(hashed)))
	fmt.Println(key.Address(bitcoinlib.COMPRESSED, true))
}
