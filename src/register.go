package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func register(ctx *gin.Context) {
	var body loginToken
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	has, err := engine.Where("account = ?", body.Account).Exist(&User{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}

	if !has {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.PW), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		_, err = engine.InsertOne(&User{Account: body.Account, PwHash: string(hashedPassword), Valid: true})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "phone number already exists",
		})
	}

}
