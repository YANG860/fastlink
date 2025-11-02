package admin

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(ctx *gin.Context) {
	// 获取所有用户
	var users []db.User
	err := db.SQLEngine.Find(&users)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 返回用户列表
	ctx.JSON(200, users)
}
