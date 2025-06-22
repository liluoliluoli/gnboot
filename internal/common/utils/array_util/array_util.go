package array_util

import (
	"hash/fnv"
	"strings"
)

// GetHashElement 根据客户端ip获取boxip
func GetHashElement(str string, clientIp string) string {
	arr := strings.Split(str, ",")
	m := arr[getIndexByHash(clientIp, len(arr))]
	return m
}

func getIndexByHash(s string, length int) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	hash := h.Sum32()
	return int(hash) % length
}
