package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"bytes"
	"encoding/hex"
	"testing"
)

func TestVersionParsing(t *testing.T) {
	version_buf := bytes.NewReader([]byte{0x01, 0x00, 0x00, 0x00, 0xff})
	result, err := bitcoinlib.NewVersionFrom(version_buf)
	if err != nil {
		t.Fatalf("Faile with error %s", err)
	}
	expected := bitcoinlib.NewVersion(1)
	if expected.Ne(*result) {
		t.Fatalf("%v is different from %v", result, expected)
	}
}

func TestTransactionScriptParsing(t *testing.T) {
	scriptHex := "6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a"
	stream, _ := hex.DecodeString(scriptHex)
	script, err := bitcoinlib.ParseScript(bytes.NewReader(stream))
	if err != nil {
		t.Fatalf("Failed decoding with error: %s", err)
	}
	deserialized := hex.EncodeToString(script.Serialize())
	if scriptHex != deserialized {
		t.Fatalf("Serialized Script is different to original: %s\n!=\n%s", scriptHex, deserialized)
	}
}

func TestInputParsing(t *testing.T) {
	expectedInput := "186f9f998a5aa6f048e51dd8419a14d8a0f1a8a2836dd734d2804fe65fa35779000000008b483045022100884d142d86652a3f47ba4746ec719bbfbd040a570b1deccbb6498c75c4ae24cb02204b9f039ff08df09cbe9f6addac960298cad530a863ea8f53982c09db8f6e381301410484ecc0d46f1918b30928fa0e4ed99f16a0fb4fde0735e7ade8416ab9fe423cc5412336376789d172787ec3457eee41c04f4938de5cc17b4a10fa336a8d752adfffffffff"
	stream, _ := hex.DecodeString(expectedInput)
	input, err := bitcoinlib.NewInputFrom(bytes.NewReader(stream))
	if err != nil {
		t.Fatalf("Failed at parsing with: %s", err)
	}
	deserialized := hex.EncodeToString(input.Serialize())
	if expectedInput != deserialized {
		t.Fatalf("Failed at inputs: %s\n!=\n%s", expectedInput, deserialized)
	}
}

func TestOutputParsing(t *testing.T) {
	expectedInput := "186f9f998a5aa6f048e51dd8419a14d8a0f1a8a2836dd734d2804fe65fa35779000000008b483045022100884d142d86652a3f47ba4746ec719bbfbd040a570b1deccbb6498c75c4ae24cb02204b9f039ff08df09cbe9f6addac960298cad530a863ea8f53982c09db8f6e381301410484ecc0d46f1918b30928fa0e4ed99f16a0fb4fde0735e7ade8416ab9fe423cc5412336376789d172787ec3457eee41c04f4938de5cc17b4a10fa336a8d752adfffffffff"
	stream, _ := hex.DecodeString(expectedInput)
	input, err := bitcoinlib.NewInputFrom(bytes.NewReader(stream))
	if err != nil {
		t.Fatalf("Failed at parsing with: %s", err)
	}
	deserialized := hex.EncodeToString(input.Serialize())
	if expectedInput != deserialized {
		t.Fatalf("Failed at inputs: %s\n!=\n%s", expectedInput, deserialized)
	}
}

func TestTransactionSerialization(t *testing.T) {
	transaction, _ := hex.DecodeString("010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600")
	parsed, _ := bitcoinlib.ParseTransaction(bytes.NewReader(transaction))
	serialized := parsed.Serialize()

	decoded := hex.EncodeToString(serialized)
	expected := hex.EncodeToString(transaction)
	if decoded != expected {
		t.Fatalf("Decoded was not the same as the original transaction:\n%s\n%s", decoded, expected)
	}
}

func TestScriptEvaluation(t *testing.T) {
  pubkey, _ := hex.DecodeString("06767695935687")
  scriptSig, _ := hex.DecodeString("0152")
  pub, err := bitcoinlib.ParsePubKey(bytes.NewReader(pubkey))
  if err != nil {
    t.Fatalf("Failed processing pub key")
  }
  sig, err := bitcoinlib.ParseScript(bytes.NewReader(scriptSig))
  if err != nil {
    t.Fatalf("Failed processing scrip sig")
  }
  combined := pub.Combine(*sig)
  if !combined.Evaluate("") {
    t.Fatalf("Failed evaluating script")
  }

}

func TestScriptEvaluationP2PK(t *testing.T) {
	z := "7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d"
	sec, _ := hex.DecodeString("04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	sig, _ := hex.DecodeString("3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601")
	pk, err := bitcoinlib.ParseFromSec(sec)
	if err != nil {
		t.Fatalf("Bad sec: %s", err)
	}
	_, err = bitcoinlib.ParseFromDer(pk, sig[:len(sig)-1])
	if err != nil {
		t.Fatalf("Bad DER: %s", err)
	}
	scriptSig := bitcoinlib.NewScript([]bitcoinlib.Operation{
		&bitcoinlib.ScriptVal{sig},
		
	})
	pubKey := bitcoinlib.NewPubkey([]bitcoinlib.Operation{
		&bitcoinlib.ScriptVal{sec},
		&bitcoinlib.OP_CHECKSIG{},
	})

	combined := pubKey.Combine(*scriptSig)
	if !combined.Evaluate(z) {
		t.Fatal("Failed evaluating Script")
	}
}

func TestTransactionInputsVsOutputs(t *testing.T) {
	tx := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	hexed, _ := hex.DecodeString(tx)
	parsed, _ := bitcoinlib.ParseTransaction(bytes.NewReader(hexed))
	if parsed.Fee(false) < 0 {
		t.Fatalf("Failed transaction: %d", parsed.Fee(false))
	}
}

func TestSigHash(t *testing.T) {
	tx, _ := bitcoinlib.FetchTransaction("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03", false, true)
	want := "27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6"
	if hex.EncodeToString(tx.SigHash(0, false)) !=  want {
		t.Fatalf("Failed sighash")
	}
}

func TestVerifiyP2PKH(t *testing.T) {
	tx, err := bitcoinlib.FetchTransaction("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03", false, true)
	if err != nil {
		t.Fatal("Failed to fetch transaction")
	}
	if !tx.Verify(false) {
		t.Fatal("Failed to verify Transaction")
	}
	tx, err = bitcoinlib.FetchTransaction("5418099cc755cb9dd3ebc6cf1a7888ad53a1a3beb5a025bce89eb1bf7f1650a2", true, true)
	if err != nil {
		t.Fatal("Failed to fetch transaction 2")
	}
	if !tx.Verify(true) {
		t.Fatal("Failed to verify Transaction 2")
	}
}