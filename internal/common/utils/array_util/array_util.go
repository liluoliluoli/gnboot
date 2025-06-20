package array_util

import (
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"hash/fnv"
)

// GetHashElement 根据客户端ip获取boxip
func GetHashElement(arr []map[string]string, clientIp string) (string, string) {
	m := arr[getIndexByHash(clientIp, len(arr))]
	if len(m) == 0 {
		return "", ""
	}
	return m[constant.Key_XiaoYaBoxIp], m[constant.Key_JellyfinBoxIp]
}

func getIndexByHash(s string, length int) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	hash := h.Sum32()
	return int(hash) % length
}
