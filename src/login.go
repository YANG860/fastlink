package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// login 处理用户登录请求。
// 步骤：
// 1. 解析请求体，获取账号和密码。
// 2. 查询数据库，验证账号是否存在。
// 3. 检查用户是否有效。
// 4. 校验密码是否正确。
// 5. 生成 JWT token 并返回给前端。
func login(ctx *gin.Context) {

	var body LoginRequest
	// 解析请求体中的 JSON 数据到 body 结构体
	err := ctx.BindJSON(&body)

	if err != nil {
		// 解析失败，返回 500 错误
		ctx.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var user User
	// 根据账号从数据库查询用户
	has, err := engine.Where("account = ?", body.Account).Get(&user)
	if err != nil {
		// 查询数据库出错，返回 500 错误
		ctx.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if !has {
		// 用户不存在，返回 401 错误
		ctx.JSON(401, Response{
			Success: false,
			Error:   "incorrect account or password",
		})
		return
	}

	if !user.Valid {
		// 用户无效，返回 401 错误
		ctx.JSON(401, Response{
			Success: false,
			Error:   "user is not valid",
		})
		return
	}

	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(body.PW))
	if err != nil {
		// 密码错误，返回 401 错误
		ctx.JSON(401, Response{
			Success: false,
			Error:   "incorrect account or password",
		})
		return
	}

	// 生成 JWT token
	tokenString, err := GenJWT(&UserToken{
		ID:      user.ID,
		Account: user.Account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		}})

	if err != nil {
		// 生成 token 失败，返回 500 错误
		ctx.JSON(500, Response{
			Success: false,
			Error:   "failed to generate token",
		})
		return
	}

	// 登录成功，返回 token
	ctx.JSON(200, LoginResponse{
		Response: Response{
			Success: true,
		},
		Token: tokenString,
	})

}
