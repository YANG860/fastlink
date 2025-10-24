package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func login(ctx *gin.Context) {

	var body LoginRequest
	err := ctx.BindJSON(&body)

	if err != nil {
		ctx.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var user User
	has, err := engine.Where("account = ?", body.Account).Get(&user)
	if err != nil {
		ctx.JSON(500, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if !has {
		ctx.JSON(401, Response{
			Success: false,
			Error:   "incorrect account or password",
		})
		return
	}

	if !user.Valid {
		ctx.JSON(401, Response{
			Success: false,
			Error:   "user is not valid",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(body.PW))
	if err != nil {
		ctx.JSON(401, Response{
			Success: false,
			Error:   "incorrect account or password",
		})
		return
	}

	tokenString, err := GenJWT(&UserToken{
		ID:      user.ID,
		Account: user.Account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		}})

	if err != nil {
		ctx.JSON(500, Response{
			Success: false,
			Error:   "failed to generate token",
		})
		return
	}

	ctx.JSON(200, LoginResponse{
		Response: Response{
			Success: true,
		},
		Token: tokenString,
	})

}
