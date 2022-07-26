package algorithm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeShorteningFunctionality(t *testing.T) {
	shr := DefaultShortener{HashSize: 10}

	hash := shr.ComputeShortening("https://minecraft.net")
	for i := 0; i < 10; i++ {
		if shr.ComputeShortening("https://minecraft.net") != hash {
			t.Fail()
		}
	}
}

func TestComputeShorteningBijection(t *testing.T) {
	shr := DefaultShortener{HashSize: 10}

	storage := make(map[string]string)

	for i := 0; i < 10000; i++ {
		url := fmt.Sprintf("https://example%d.com", i+1)
		hash := shr.ComputeShortening(url)

		if value, ok := storage[hash]; ok {
			fmt.Printf("Collision found: %s and %s\n", value, url)
			t.Fail()
		}

		storage[hash] = url
	}
}

func TestStubShortener(t *testing.T) {
	shr := StubShortener{}

	assert.NotEqual(t, "", shr.ComputeShortening("test-url"))
}
