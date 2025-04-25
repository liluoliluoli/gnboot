package security_util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func SignBySha256(bizParamStr string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(bizParamStr))
	hashValue := hash.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(hashValue)
	return sign, nil
}

func SignByHMACSha256(message, key string) string {
	keyBytes := []byte(key)
	hash := hmac.New(sha256.New, keyBytes)
	hash.Write([]byte(message))
	hashValue := hash.Sum(nil)
	hexHashValue := hex.EncodeToString(hashValue)
	return hexHashValue
}
