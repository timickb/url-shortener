package algorithm

import (
	"crypto/md5"
	"github.com/speps/go-hashids/v2"
	"strconv"
)

const (
	lower    = "abcdefghijklmnopqrstuvwxyz"
	upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers  = "0123456789"
	special  = "_"
	alphabet = lower + upper + numbers + special
	hashSize = 10
)

func ComputeShortening(str string) string {
	var salt int
	for i := 0; i < len(str); i++ {
		salt += int(str[i])
	}
	md5hash := md5.Sum([]byte(str))

	md5sumInt := make([]int, len(md5hash))
	for i := 0; i < len(md5hash); i++ {
		md5sumInt[i] = int(md5hash[i])
	}

	hd := hashids.NewData()
	hd.MinLength = hashSize
	hd.Alphabet = alphabet
	hd.Salt = strconv.Itoa(salt)
	h, _ := hashids.NewWithData(hd)
	encoded, _ := h.Encode(md5sumInt)

	return encoded[:hashSize]
}
