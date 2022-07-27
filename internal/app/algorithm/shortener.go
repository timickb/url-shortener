package algorithm

import (
	"crypto/md5"
	"strconv"

	"github.com/speps/go-hashids/v2"
)

const (
	lower    = "abcdefghijklmnopqrstuvwxyz"
	upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers  = "0123456789"
	special  = "_"
	alphabet = lower + upper + numbers + special
)

type Shortener interface {
	ComputeShortening(url string) string
}

type DefaultShortener struct {
	HashSize int
}

func (s DefaultShortener) ComputeShortening(url string) string {
	var salt int
	for i := 0; i < len(url); i++ {
		salt += int(url[i])
	}
	md5hash := md5.Sum([]byte(url))

	md5sumInt := make([]int, len(md5hash))
	for i := 0; i < len(md5hash); i++ {
		md5sumInt[i] = int(md5hash[i])
	}

	hd := hashids.NewData()
	hd.MinLength = s.HashSize
	hd.Alphabet = alphabet
	hd.Salt = strconv.Itoa(salt)
	h, _ := hashids.NewWithData(hd)
	encoded, _ := h.Encode(md5sumInt)

	return encoded[:s.HashSize]
}
