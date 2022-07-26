package algorithm

import (
	"fmt"
	"testing"
)

func TestComputeShorteningFunctionality(t *testing.T) {
	t.Helper()

	hash := ComputeShortening("https://minecraft.net")
	for i := 0; i < 10; i++ {
		if ComputeShortening("https://minecraft.net") != hash {
			t.Fail()
		}
	}
}

func TestComputeShorteningBijection(t *testing.T) {
	t.Helper()

	storage := make(map[string]string)

	for i := 0; i < 10000; i++ {
		url := fmt.Sprintf("https://example%d.com", i+1)
		hash := ComputeShortening(url)

		if value, ok := storage[hash]; ok {
			fmt.Printf("Collision found: %s and %s\n", value, url)
			t.Fail()
		}

		storage[hash] = url
	}
}
