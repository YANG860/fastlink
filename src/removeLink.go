package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

func removeLink(ctx *gin.Context, url string) {
	var body RemoveLinkRequest
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, RemoveLinkResponse{Response: Response{Success: false, Error: "Invalid request"}})
		return
	}

	userToken, err := ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, RemoveLinkResponse{Response: Response{Success: false, Error: "Invalid token"}})
		return
	}

	var link Link
	has, err := engine.Where("short_url=?", url).Get(&link)
	if err != nil {
		ctx.JSON(500, RemoveLinkResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}
	if !has || link.UserID != userToken.ID {
		ctx.JSON(404, RemoveLinkResponse{Response: Response{Success: false, Error: "Link not found"}})
		return
	}

	_, err = engine.ID(link.ID).Update(&Link{ExpireAt: time.Now().Add(-time.Hour)})
	if err != nil {
		ctx.JSON(500, RemoveLinkResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}

	ctx.JSON(200, RemoveLinkResponse{Response: Response{Success: true}})
}
