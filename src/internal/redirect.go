package internal

import (
	"fastlink/db"
	"fastlink/models"
	"time"

	"github.com/gin-gonic/gin"
)

// 添加redis缓存

// Redirect 根据短链重定向到原始链接，并增加点击数
func Redirect(ctx *gin.Context, shortUrl string) {
	var err error
	var link *db.Link
	link, err = db.GetLinkFromCache(shortUrl)

	if err != nil {
		// 缓存未命中，回表查询
		_, err := db.SQLEngine.Where("short_url=?", shortUrl).Get(link)
		if err != nil {
			ctx.JSON(500, models.DatabaseError)
			return
		}
		if link.ExpireAt.Before(time.Now()) {
			ctx.JSON(404, models.NotFoundError)
			db.SetLinkToCache(&db.Link{ShortUrl: shortUrl, ExpireAt: link.ExpireAt})
			return
		}

		//创建者是否有效
		var user db.User
		_, err = db.SQLEngine.ID(link.UserID).Get(&user)

		if err != nil {
			ctx.JSON(500, models.DatabaseError)
			return
		}
		if user.Valid == false {
			ctx.JSON(404, models.NotFoundError)
			db.SetLinkToCache(&db.Link{ShortUrl: shortUrl, ExpireAt: link.ExpireAt})
			return
		}

		err = db.SetLinkToCache(link)
		if err != nil {
			ctx.JSON(500, models.InternalServerError)
			return
		}
	}

	if link.SourceUrl==""{
		ctx.JSON(404, models.NotFoundError)
		return
	}

	// 重定向到原始链接
	ctx.Redirect(302, link.SourceUrl)

	// 异步增加点击数
	
}
