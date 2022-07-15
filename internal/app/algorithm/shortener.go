package algorithm

import (
	"strings"
)

const (
	lower    = "abcdefghijklmnopqrstuvwxyz"
	upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers  = "0123456789"
	special  = "_"
	alphabet = lower + upper + numbers + special
	alpLen   = uint(len(alphabet))
)

func ComputeHash(str string, resultSize uint) string {
	strLen := uint(len(str))
	data := strings.Clone(str)
	addition := resultSize*2 - (strLen % resultSize)
	blockSize := (strLen + addition) / resultSize

	var chrSum uint
	for i := uint(0); i < strLen; i++ {
		chrSum += uint(str[i])
	}

	// We want len(data) % resultSize = 0; adding symbols to the end
	for i := uint(1); i <= addition; i++ {
		// added symbols depends on symbols in source string
		data += string(alphabet[(chrSum*i+uint(str[i%strLen]))%alpLen])
	}

	// divide source string on blocks of size ceil(strLen / resultSize)
	var result string
	for block := uint(0); block < resultSize; block++ {
		// each block maps to one symbol in result string
		// each result symbol depends on source string size and block char sum
		var blockSum uint
		for i := block * blockSize; i < block*blockSize+blockSize; i++ {
			blockSum += uint(data[i])
		}
		dependency := ((strLen * (blockSum + resultSize)) % resultSize) * (blockSize + 1)
		result += string(alphabet[dependency%alpLen])
	}

	return result
}
