package auth

import (
	"fastlink/db"
	"fastlink/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// register 用户注册处理函数
func Register(ctx *gin.Context) {
	var body models.RegisterRequest
	// 解析请求体
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	// 检查账号是否已存在
	has, err := db.SQLEngine.Where("account = ?", body.Account).Exist(&db.User{})

	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	if has {
		ctx.JSON(400, models.AlreadyExistsError)
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.PW), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, models.InternalServerError)
		return
	}

	// 插入新用户
	_, err = db.SQLEngine.InsertOne(&db.User{Account: body.Account, PwHash: string(hashedPassword), Valid: true})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	ctx.JSON(200, models.Success)

}
