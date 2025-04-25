package jwtutil

import (
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func ParseJWT(tokenString string, secretKey string) (jwt.MapClaims, bool, error) {
	// 兼容Java版本，先解密base64，得到 []byte
	secret, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, false, err
	}
	// 解析jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect: unexpected signing method: token.Header["alg"]
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || jwt.SigningMethodHS512.Alg() != token.Header["alg"] {
			return nil, errors.New("invalid token")
		}
		// decodeSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})
	if err != nil {
		return nil, false, err
	}
	// Check if the token is valid and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true, nil
	}
	return nil, false, errors.New("invalid token")
}

func GenerateJwt(claims jwt.MapClaims, secretKey string) (string, error) {
	//// 兼容Java版本，先解密base64，得到 []byte
	//secret, err := base64.StdEncoding.DecodeString(secretKey)
	//if err != nil {
	//	return "", err
	//}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateUserToken(claims *UserClaims, secretKey string) (string, error) {
	if claims == nil {
		return "", nil
	}
	mapClaims := jwt.MapClaims{
		"userName":  claims.UserName,
		"timestamp": time.Now().UnixMilli(),
		"iss":       "yunvd",
	}
	return GenerateJwt(mapClaims, secretKey)
}

type UserClaims struct {
	UserName string `json:"userName"`
}
