package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// register 用户注册处理函数
func register(ctx *gin.Context) {
	var body LoginRequest
	// 解析请求体
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(500, Response{
			Success: false,
			Error:   "Invalid request",
		})
		return
	}

	// 检查账号是否已存在
	has, err := engine.Where("account = ?", body.Account).Exist(&User{})

	if err != nil {
		ctx.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return

	}

	if has {
		ctx.JSON(500, Response{
			Success: false,
			Error:   "account already exists",
		})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.PW), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, Response{
			Success: false,
			Error:   "failed to hash password",
		})
		return
	}

	// 插入新用户
	_, err = engine.InsertOne(&User{Account: body.Account, PwHash: string(hashedPassword), Valid: true})
	if err != nil {
		ctx.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(200, Response{
		Success: true,
	})

}
