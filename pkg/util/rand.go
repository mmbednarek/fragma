package util

import (
	"math/rand"
)

func String(n int) string {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = byte('a' + rand.Int31n('z'-'a'))
	}
	return string(result)
}
