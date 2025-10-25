package main

import (
	"github.com/gin-gonic/gin"
)

// getLink 获取短链详细信息，需验证 token 和权限
func getLink(ctx *gin.Context, short string) {

	var body GetLinkRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, GetLinkResponse{Response: Response{Success: false, Error: "Invalid request"}})
		return
	}

	var link Link
	// 查询短链
	has, err := engine.Where("short_url=?", short).Get(&link)
	if err != nil {
		ctx.JSON(500, GetLinkResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}
	if !has {
		ctx.JSON(404, GetLinkResponse{Response: Response{Success: false, Error: "Link not found"}})
		return
	}
	// 校验 token
	token, err := ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, GetLinkResponse{Response: Response{Success: false, Error: "Invalid token"}})
		return
	}
	// 权限校验
	if link.UserID != token.ID {
		ctx.JSON(403, GetLinkResponse{Response: Response{Success: false, Error: "Forbidden"}})
		return
	}
	// 返回短链信息
	ctx.JSON(200, GetLinkResponse{
		Response:   Response{Success: true},
		SourceUrl:  link.SourceUrl,
		ShortUrl:   link.ShortUrl,
		ClickCount: link.ClickCount,
		CreatedAt:  link.CreatedAt,
		ExpireAt:   link.ExpireAt,
	})
}
