package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"encoding/hex"
	"testing"
)

func TestMerkleParent(t *testing.T) {
	hash0, _ := hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	hash1, _ := hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")

	expected := "8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd"

	obtained := hex.EncodeToString(bitcoinlib.MerkleParent(hash0, hash1))

	if obtained != expected {
		t.Fatalf("Failed creating parent: %s != %s", expected, obtained)
	}
}

func TestMerkleParentLevel(t *testing.T) {
	children := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
		"2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908",
	}
	expectedParents := []string{
		"8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd",
		"7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800",
		"ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7",
		"68b3e2ab8182dfd646f13fdf01c335cf32476482d963f5cd94e934e6b3401069",
		"43e7274e77fbe8e5a42a8fb58f7decdb04d521f319f332d88e6b06f8e6c09e27",
		"1796cd3ca4fef00236e07b723d3ed88e1ac433acaaa21da64c4b33c946cf3d10",
	}
	convertedChildren := make([][]byte, len(children))

	for i, child := range children {
		convertedChildren[i], _ = hex.DecodeString(child)
	}

	obtainedParents := bitcoinlib.MerkleParentLevel(convertedChildren)

	for i, obtained := range obtainedParents {
		parent := hex.EncodeToString(obtained)
		if parent != expectedParents[i] {
			t.Fatalf("Failed at index %d: %s vs. %s", i, expectedParents[i], parent)
		}
	}
}

func TestMerkleRoot(t *testing.T) {
	children := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
		"2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908",
		"b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0",
	}
	expectedRoot := "acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6"
	convertedChildren := make([][]byte, len(children))

	for i, child := range children {
		convertedChildren[i], _ = hex.DecodeString(child)
	}

	obtainedRoot := hex.EncodeToString(bitcoinlib.MerkleRoot(convertedChildren))

	if obtainedRoot != expectedRoot {
		t.Fatalf("Expected merkle root %s but got %s", expectedRoot, obtainedRoot)
	}
}

func TestMerkleBlock(t *testing.T) {
	blockHex := "00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635"
	merkleBytes, _ := hex.DecodeString(blockHex)
	merkle := bitcoinlib.NewMerkleBlockMessage()
	_, err := merkle.Parse(merkleBytes)
	if err != nil {
		t.Fatalf("Failed at parsing MerkleBlock: %s", err)
	}
	decoded := merkle.Serialize()
	if hex.EncodeToString(decoded) != blockHex {
		t.Fatalf("Failed reserializing Merkle Block: %s vs %s", hex.EncodeToString(decoded), blockHex)
	}
}
