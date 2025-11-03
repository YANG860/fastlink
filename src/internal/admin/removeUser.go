package admin

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

func RemoveUser(ctx *gin.Context) {

	var body models.RemoveUserAdminRequest
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

	
	// 逻辑删除用户
	_, err = db.SQLEngine.ID(body.UserID).Update(&db.User{Valid: true})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 获取该用户的所有短链
	var links []db.Link
	err = db.SQLEngine.Where("user_id = ?", body.UserID).Find(&links)
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
