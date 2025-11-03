package admin

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

func RemoveUser(ctx *gin.Context, userID int) {
	// 逻辑删除用户
	_, err := db.SQLEngine.ID(userID).Update(&db.User{Valid: true})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 获取该用户的所有短链
	var links []db.Link
	err = db.SQLEngine.Where("user_id = ?", userID).Find(&links)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 删除缓存
	for _, link := range links {
		err = db.SetLinkToCache(&db.Link{ShortUrl: link.ShortUrl, ExpireAt: link.ExpireAt})
		if err != nil {
			ctx.JSON(500, models.InternalServerError)
			return
		}
	}

	ctx.JSON(200, models.Success)
}
