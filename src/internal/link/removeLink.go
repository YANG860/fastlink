package link

import (
	"fastlink/auth"
	"fastlink/db"
	"fastlink/internal/admin"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

// removeLink 逻辑删除短链（提前过期），需验证 token 和权限
func RemoveLink(ctx *gin.Context, url string) {
	var body models.RemoveLinkRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	// 校验 token
	userToken, err := auth.ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	var link db.Link
	// 查询短链
	has, err := db.SQLEngine.Where("short_url=?", url).Get(&link)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}
	// 权限校验
	if !has || link.UserID != userToken.ID {
		ctx.JSON(404, models.NotFoundError)
		return
	}

	admin.RemoveLink(ctx, url)
}
