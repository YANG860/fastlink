package link

import (
	"fastlink/auth"
	"fastlink/db"
	"fastlink/models"
	"time"

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

	_, err = db.SQLEngine.Where("short_url = ?", url).Update(&db.Link{ExpireAt: time.Now().Add(-time.Minute)})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 删除缓存
	err = db.SetLinkToCache(&db.Link{ShortUrl: url, ExpireAt: time.Now().Add(-time.Minute)})
	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	ctx.JSON(200, models.Success)
}
