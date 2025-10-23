package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func register(ctx *gin.Context) {
	var body LoginRequest
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request",
		})
		return
	}

	has, err := engine.Where("account = ?", body.Account).Exist(&User{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return

	}

	if has {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "account already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.PW), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "failed to hash password",
		})
		return
	}

	_, err = engine.InsertOne(&User{Account: body.Account, PwHash: string(hashedPassword), Valid: true})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
	})

}
