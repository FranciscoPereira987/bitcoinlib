package bitcoinlib

import "io"

type Script struct {}

func ParseScript(from io.Reader) []byte {
	last_four := make([]byte, 4)
	new_four := make([]byte, 4)
	val, _ := from.Read(new_four)
	for  val >= 4 {
		copy(last_four, new_four)
	}
	copy(last_four[3-val:], new_four[:3-val])
	return last_four
}