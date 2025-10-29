package internal

import (
	"fastlink/db"
	"fastlink/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 添加redis缓存

// Redirect 根据短链重定向到原始链接，并增加点击数
func Redirect(ctx *gin.Context, shortUrl string) {

	var source string
	var linkID int
	var err error
	source, err = db.RedisClient.Get(db.Ctx, "short_url:"+shortUrl+":source").Result()
	//TODO: 重构缓存处理 点击计数  缓存异步写回

	if err != nil {
		// 缓存未命中，回表查询
		var link db.Link
		_, err := db.Engine.Where("short_url=?", shortUrl).Get(&link)
		if err != nil {
			ctx.JSON(500, models.DatabaseError)
			return
		}

		// 判断短链是否过期 返回404，并将空值写入缓存，防止缓存穿透
		if link.ExpireAt.Before(time.Now()) {
			ctx.JSON(404, models.NotFoundError)
			db.RedisClient.Set(db.Ctx, "short_url:"+shortUrl+":source", "", 30*time.Minute)
			return
		}

		//创建者是否有效
		var user db.User
		_, err = db.Engine.ID(link.UserID).Get(&user)
		if err != nil || !user.Valid {
			ctx.JSON(404, models.NotFoundError)
			db.RedisClient.Set(db.Ctx, "short_url:"+shortUrl+":source", "", 30*time.Minute)
			return
		}

		source = link.SourceUrl
		err = db.RedisClient.Set(db.Ctx, "short_url:"+shortUrl+":source", source, 0).Err()
		if err != nil {
			ctx.JSON(500, models.InternalServerError)
			return
		}
		err = db.RedisClient.Set(db.Ctx, "short_url:"+shortUrl+":link_id", link.ID, 0).Err()
		if err != nil {
			ctx.JSON(500, models.InternalServerError)
			return
		}
	}
	// 重定向到原始链接
	if source == "" {
		ctx.JSON(404, models.NotFoundError)
		return
	}
	ctx.Redirect(301, source)

	// 增加点击数
	linkID_str, err := db.RedisClient.Get(db.Ctx, "short_url:"+shortUrl+":link_id").Result()
	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	linkID, err = strconv.Atoi(linkID_str)
	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	var link db.Link
	_, err = db.Engine.ID(linkID).Get(&link)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}
	db.Engine.ID(linkID).Update(&db.Link{
		ClickCount: link.ClickCount + 1,
	})

}
