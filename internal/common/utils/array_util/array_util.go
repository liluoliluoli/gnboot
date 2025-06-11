package array_util

import (
	"crypto/rand"
	"math/big"
)

func GetRandomElement(arr []string) string {
	n := big.NewInt(int64(len(arr)))
	index, _ := rand.Int(rand.Reader, n)
	return arr[index.Int64()]
}
