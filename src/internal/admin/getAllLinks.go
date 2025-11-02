package admin

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

func GetAllLinks(ctx *gin.Context) {
	// 获取所有短链
	var links []db.Link
	err := db.SQLEngine.Find(&links)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 返回短链列表
	ctx.JSON(200, links)
}
