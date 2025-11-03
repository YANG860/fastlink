package admin

import (
	"fastlink/db"
	"fastlink/models"
	"time"

	"github.com/gin-gonic/gin"
)

func RemoveLink(ctx *gin.Context) {

	var body models.RemoveLinkAdminRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	ok, err := db.AuthenticateAdmin(body.Admin.Account, body.Admin.PW)
	if !ok {
		ctx.JSON(403, models.ForbiddenError)
		return
	}
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 逻辑删除短链（提前过期）
	_, err = db.SQLEngine.Where("short_url = ?", body.ShortUrl).Update(&db.Link{ExpireAt: time.Now().Add(-time.Minute)})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}


	// 删除缓存
	err = db.SetLinkToCache(&db.Link{ShortUrl: body.ShortUrl, ExpireAt: time.Now().Add(-time.Minute)})
	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	ctx.JSON(200, models.Success)
}
