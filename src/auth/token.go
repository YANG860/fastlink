package auth

import "github.com/golang-jwt/jwt/v5"

// JWT 用户信息结构体
type UserToken struct {
	jwt.RegisteredClaims
	ID      int    `json:"id"`
	Account string `json:"account"`
}
