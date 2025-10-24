package main

import (
	"github.com/gin-gonic/gin"
)

func getLink(ctx *gin.Context, short string) {

	var body GetLinkRequest
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, GetLinkResponse{Response: Response{Success: false, Error: "Invalid request"}})
		return
	}


	var link Link
	has, err := engine.Where("short_url=?", short).Get(&link)
	if err != nil {
		ctx.JSON(500, GetLinkResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}
	if !has {
		ctx.JSON(404, GetLinkResponse{Response: Response{Success: false, Error: "Link not found"}})
		return
	}
	token, err := ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, GetLinkResponse{Response: Response{Success: false, Error: "Invalid token"}})
		return
	}
	if link.UserID != token.ID {
		ctx.JSON(403, GetLinkResponse{Response: Response{Success: false, Error: "Forbidden"}})
		return
	}
	ctx.JSON(200, GetLinkResponse{
		Response:    Response{Success: true},
		SourceUrl:   link.SourceUrl,
		ShortUrl:    link.ShortUrl,
		ClickCount:  link.ClickCount,
		CreatedAt:   link.CreatedAt,
		ExpireAt:    link.ExpireAt,
	})
}