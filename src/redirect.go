package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func redirect(ctx *gin.Context, short string) {

	var link Link
	has, err := engine.Where("short_url=?", short).Get(&link)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if !has {
		ctx.JSON(404, gin.H{"error": "Not found"})
		return
	}

	ctx.Redirect(http.StatusPermanentRedirect, link.SourceUrl)

}
