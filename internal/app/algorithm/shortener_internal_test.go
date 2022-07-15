package algorithm

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_ComputeHash(t *testing.T) {
	t.Helper()

	url1 := "http://example.com"
	url2 := "https://google.com"

	hash1 := ComputeHash(url1)
	hash2 := ComputeHash(url2)

	log.Println(hash1)
	log.Println(hash2)

	assert.NotEqual(t, hash1, hash2)

}
