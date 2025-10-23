// 登录相关主程序
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func login(ctx *gin.Context) {

	var body LoginRequest
	err := ctx.BindJSON(&body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var user User
	has, err := engine.Where("account = ?", body.Account).Get(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if !has {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "incorrect account or password",
		})
		return
	}

	if !user.Valid {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "user is not valid",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(body.PW))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "incorrect account or password",
		})
		return
	}

	tokenString, err := GenJWT(&UserToken{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		}})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Response: Response{
			Success: true,
		},
		Token: tokenString,
	})

}
