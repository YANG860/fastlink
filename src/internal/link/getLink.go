package link

import (
	"fastlink/auth"
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

// getLink 获取短链详细信息，需验证 token 和权限
func GetLink(ctx *gin.Context, short string) {

	var body models.GetLinkRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	var link db.Link
	// 查询短链
	has, err := db.SQLEngine.Where("short_url=?", short).Get(&link)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}
	if !has {
		ctx.JSON(404, models.NotFoundError)
		return
	}
	// 校验 token
	token, err := auth.ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, models.InvalidTokenError)
		return
	}
	// 权限校验
	if link.UserID != token.ID {
		ctx.JSON(403, models.ForbiddenError)
		return
	}
	// 返回短链信息
	ctx.JSON(200, models.GetLinkResponse{
		Response:   models.Success,
		SourceUrl:  link.SourceUrl,
		ShortUrl:   link.ShortUrl,
		ClickCount: link.ClickCount,
		CreatedAt:  link.CreatedAt,
		ExpireAt:   link.ExpireAt,
	})
}
