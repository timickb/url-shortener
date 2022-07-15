package algorithm

import (
	"fmt"
	"testing"
)

func TestComputeShortening_functionality(t *testing.T) {
	t.Helper()

	hash := ComputeShortening("https://minecraft.net")
	for i := 0; i < 10; i++ {
		if ComputeShortening("https://minecraft.net") != hash {
			t.Fail()
		}
	}
}

func TestComputeShortening_bijection(t *testing.T) {
	t.Helper()

	storage := make(map[string]string)

	for i := 0; i < 10000; i++ {
		url := fmt.Sprintf("https://example%d.com", i+1)
		hash := ComputeShortening(url)

		fmt.Println(hash)

		if value, ok := storage[hash]; ok {
			fmt.Printf("Collision found: %s and %s\n", value, url)
			t.Fail()
		}

		storage[hash] = url
	}
}
