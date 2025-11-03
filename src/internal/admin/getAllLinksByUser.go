package admin

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

func GetAllLinksByUser(ctx *gin.Context) {

	var body models.GetAllLinksByUserAdminRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	ok, err := db.AuthenticateAdmin(body.Admin.Account, body.Admin.PW)
	if !ok {
		ctx.JSON(403, models.ForbiddenError)
		return
	}

	// 获取所有短链
	var links []db.Link
	err = db.SQLEngine.Where("user_id = ?", body.UserID).Find(&links)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 返回短链列表
	ctx.JSON(200, models.GetAllLinksByUserAdminResponse{
		Response: models.Success,
		Links:    links,
	})
}
