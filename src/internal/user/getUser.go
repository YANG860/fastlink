package user

import (
	"fastlink/auth"
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
)

// getUser 获取用户信息，需验证 token 和权限
func GetUser(ctx *gin.Context, account string) {

	var body models.GetUserRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	// 校验 token
	token, err := auth.ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	// 权限校验
	if token.Account != account {
		ctx.JSON(403, models.ForbiddenError)
		return
	}

	var user db.User
	// 查询用户
	has, err := db.SQLEngine.ID(token.ID).Get(&user)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}
	if !has {
		ctx.JSON(404, models.NotFoundError)
		return
	}

	// 返回用户信息
	ctx.JSON(200, models.GetUserResponse{
		Response:        models.Success,
		Account:         user.Account,
		RegisteredAt:    user.RegisteredAt,
		LinkCount:       user.LinkCount,
		LatestCreatedAt: user.LatestCreatedAt,
	})
}
