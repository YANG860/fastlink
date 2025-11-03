package admin

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(ctx *gin.Context) {

	var body models.GetAllUsersAdminRequest
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

	// 获取所有用户
	var users []db.User
	err = db.SQLEngine.Find(&users)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 返回用户列表
	ctx.JSON(200, models.GetAllUsersAdminResponse{
		Response: models.Success,
		Users:    users,
	})
}
