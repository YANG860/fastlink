package user

import (
	"fastlink/auth"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

//need to delete links

// removeUser 注销用户（逻辑删除），需验证 token 和权限
func RemoveUser(ctx *gin.Context, account string) {

	var body models.RemoveUserRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	// 校验 token
	userToken, err := auth.ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	// 权限校验
	if userToken.Account != account {
		ctx.JSON(403, models.ForbiddenError)
		return
	}

	// 逻辑删除用户
	_, err = models.Engine.ID(userToken.ID).Cols("valid").Update(&models.User{Valid: false})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}
	ctx.JSON(200, models.RemoveUserResponse{Response: models.Success})
}
