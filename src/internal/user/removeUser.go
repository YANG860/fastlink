package user

import (
	"fastlink/auth"
	"fastlink/db"
	"fastlink/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// removeUser 注销用户（逻辑删除），需验证 token 和权限
// 并发一致性

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
	_, err = db.SQLEngine.ID(userToken.ID).Cols("valid").Update(&db.User{Valid: false})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	//TODO: 触发相关资源的删除  封装进一个函数
	//  删除用户相关的短链

	//  redis 缓存 version++
	var tokenID_str = strconv.Itoa(userToken.ID)

	err = db.RedisClient.Set(db.Ctx, "user:"+tokenID_str+":token_version", userToken.Version+1, 0).Err()
	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	//  数据库version++
	_, err = db.SQLEngine.ID(userToken.ID).Cols("token_version").Update(&db.User{TokenVersion: userToken.Version + 1})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	ctx.JSON(200, models.RemoveUserResponse{Response: models.Success})
}
