package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

func redirect(ctx *gin.Context, short string) {

	var link Link
	has, err := engine.Where("short_url=?", short).Get(&link)
	if err != nil {
		ctx.JSON(500, Response{Success: false, Error: "Database error"})
		return
	}
	if link.ExpireAt.Before(time.Now()) {
		ctx.JSON(404, Response{Success: false, Error: "Not found"})
		return
	}

	_, err = engine.ID(link.ID).Update(&Link{ClickCount: link.ClickCount + 1})

	if err != nil {
		ctx.JSON(500, Response{Success: false, Error: "Internal server error"})
		return
	}

	if !has {
		ctx.JSON(404, Response{Success: false, Error: "Not found"})
		return
	}

	ctx.Redirect(301, link.SourceUrl)

}
