package admin

import (
	"fastlink/db"
	"fastlink/models"
	"time"

	"github.com/gin-gonic/gin"
)

func RemoveLink(ctx *gin.Context, shortUrl string) {
	// 逻辑删除短链（提前过期）
	_, err := db.SQLEngine.Where("short_url = ?", shortUrl).Update(&db.Link{ExpireAt: time.Now().Add(-time.Minute)})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 删除缓存
	err = db.SetLinkToCache(&db.Link{ShortUrl: shortUrl, ExpireAt: time.Now().Add(-time.Minute)})
	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	ctx.JSON(200, models.Success)
}
