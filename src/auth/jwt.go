package auth

import (
	"fastlink/db"

	"fmt"

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

	// 解析 token
	// 首先检查redis 缓存
	// 回表查询是否是最新版本 版本号加入缓存  30min过期
	// 如果不是最新版本 则提示重新登录

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

	// 从 Redis 获取用户的 token 版本号

	var version int
	version, err = db.GetVersionFromCache(token.ID)

	if err != nil {

		// 缓存未命中，回表查询
		var user db.User
		_, err := db.SQLEngine.ID(token.ID).Get(&user)
		if err != nil {
			return nil, err
		}
		version = user.TokenVersion

		// 更新缓存
		err = db.SetVersionToCache(token.ID, version)
		if err != nil {
			return nil, err
		}
	}

	if version != token.Version {
		return nil, fmt.Errorf("token expired")
	}

	return token, nil

}
