// 登录相关主程序
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 登录接口，校验用户信息并返回JWT token
func login(ctx *gin.Context) {

	// 解析请求体，获取手机号和密码
	var body loginToken
	err := ctx.BindJSON(&body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 查询数据库，查找用户信息
	var user User
	has, err := engine.Where("account = ?", body.Account).Get(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 用户不存在
	if !has {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect account or password"})
		return
	}

	// 校验密码（bcrypt）
	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(body.PW))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect account or password"})
		return
	}

	// 生成JWT token，包含手机号和过期时间
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// 使用密钥签名token
	secretKey := "key"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	// 返回token给前端
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   tokenString,
	})

}
