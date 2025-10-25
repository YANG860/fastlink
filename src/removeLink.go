package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

// removeLink 逻辑删除短链（提前过期），需验证 token 和权限
func removeLink(ctx *gin.Context, url string) {
	var body RemoveLinkRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, RemoveLinkResponse{Response: Response{Success: false, Error: "Invalid request"}})
		return
	}

	// 校验 token
	userToken, err := ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, RemoveLinkResponse{Response: Response{Success: false, Error: "Invalid token"}})
		return
	}

	var link Link
	// 查询短链
	has, err := engine.Where("short_url=?", url).Get(&link)
	if err != nil {
		ctx.JSON(500, RemoveLinkResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}
	// 权限校验
	if !has || link.UserID != userToken.ID {
		ctx.JSON(404, RemoveLinkResponse{Response: Response{Success: false, Error: "Link not found"}})
		return
	}

	// 逻辑删除（设置过期时间为过去）
	_, err = engine.ID(link.ID).Update(&Link{ExpireAt: time.Now().Add(-time.Hour)})
	if err != nil {
		ctx.JSON(500, RemoveLinkResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}

	ctx.JSON(200, RemoveLinkResponse{Response: Response{Success: true}})
}
