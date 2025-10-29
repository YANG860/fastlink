package internal

import (
	"fastlink/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Redirect 根据短链重定向到原始链接，并增加点击数
func Redirect(ctx *gin.Context, short string) {

	var link models.Link
	// 查询短链
	has, err := models.Engine.Where("short_url=?", short).Get(&link)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)

		return
	}
	if link.ExpireAt.Before(time.Now()) {
		ctx.JSON(404, models.NotFoundError)
		return
	}

	// 增加点击数
	_, err = models.Engine.ID(link.ID).Update(&models.Link{ClickCount: link.ClickCount + 1})

	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	if !has {
		ctx.JSON(404, models.NotFoundError)
		return
	}

	// 重定向到原始链接
	ctx.Redirect(301, link.SourceUrl)

}
