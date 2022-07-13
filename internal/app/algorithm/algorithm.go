package algorithm

import "fmt"

func ComputeHash(str string) string {
	return fmt.Sprintf("abacaba%d", len(str))
}
