package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"testing"
)

func TestConversionFromHexValues(t *testing.T) {
	values := []string{
		"0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d",
		"0xc7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6",
    "0xeff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c",
	}
	expectedResults := []string{
		
		"9MA8fRQrT4u8Zj8ZRd6MAiiyaxb2Y1CMpvVkHQu5hVM6",
		"EQJsjkd6JaGwxrjEhfeqPenqHwrBmPQZjJGNSCHBkcF7",
    	"4fE3H2E6XMp4SsxtwinF7w9a34ooUrwWe4WsW1458Pd",
	}

	for index, value := range values {
    encoded := bitcoinlib.IntoBase58(value[2:])
    expected := expectedResults[index]
    if encoded != expected {
      t.Fatalf("Failed at index %d\nExpected => %s\nActual => %s\n", index, expected, encoded)
    }
	}
}

func TestConversionIntoHexValues(t *testing.T) {
	expectedResults := []string{
		"000000000000006f507b27411ccf7f16f10297de6cef3f291623eddf192d9df2",
	}
	values := []string{
		"mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf",
	}

	for index, value := range values {
    encoded := bitcoinlib.FromBase58(value)
    expected := expectedResults[index]
    if encoded != expected {
      t.Fatalf("Failed at index %d\nExpected => %s\nActual => %s\n", index, expected, encoded)
    }
	}
}



func TestIntoHexValuesAdress(t *testing.T) {
	expectedResults := []string{
		"507b27411ccf7f16f10297de6cef3f291623eddf",
	}
	values := []string{
		"mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf",
	}

	for index, value := range values {
    encoded := bitcoinlib.FromBase58Address(value)
    expected := expectedResults[index]
    if encoded != expected {
      t.Fatalf("Failed at index %d\nExpected => %s\nActual => %s\n", index, expected, encoded)
    }
	}
}
