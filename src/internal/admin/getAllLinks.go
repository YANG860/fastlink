package admin

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

func GetAllLinks(ctx *gin.Context) {
	// 获取所有短链
	var body models.GetAllLinksAdminRequest
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

	var links []db.Link
	err = db.SQLEngine.Find(&links)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 返回短链列表
	ctx.JSON(200, models.GetAllLinksAdminResponse{
		Response: models.Success,
		Links:    links,
	})
}
