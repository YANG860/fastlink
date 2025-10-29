package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWT 密钥
var secretKey string = "key"

// 生成 JWT token
func GenJWT(user *UserToken) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

// 解析 JWT token
func ParseJWT(tokenString string) (*UserToken, error) {

	claims := &UserToken{}
	rawToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !rawToken.Valid {
		return nil, err
	}

	token, ok := rawToken.Claims.(*UserToken)
	if !ok {
		return nil, err
	}
	return token, nil

}
