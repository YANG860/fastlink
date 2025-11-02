package auth

import (
	"fastlink/config"
	"fastlink/db"
	"fastlink/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// login 处理用户登录请求。
// 步骤：
// 1. 解析请求体，获取账号和密码。
// 2. 查询数据库，验证账号是否存在。
// 3. 检查用户是否有效。
// 4. 校验密码是否正确。
// 5. 生成 JWT token 并返回给前端。

func Login(ctx *gin.Context) {

	var body models.LoginRequest
	// 解析请求体中的 JSON 数据到 body 结构体
	err := ctx.BindJSON(&body)

	if err != nil {
		// 解析失败，返回 500 错误
		ctx.JSON(500, models.InternalServerError)
		return
	}

	var user db.User
	// 根据账号从数据库查询用户
	has, err := db.SQLEngine.Where("account = ?", body.Account).Get(&user)
	if err != nil {
		// 查询数据库出错，返回 500 错误
		ctx.JSON(500, models.DatabaseError)
		return
	}

	if !has {
		// 用户不存在，返回 401 错误
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	if !user.Valid {
		// 用户无效，返回 401 错误
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(body.PW))
	if err != nil {
		// 密码错误，返回 401 错误
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	// 生成 JWT token
	tokenString, err := GenJWT(&UserToken{
		ID:      user.ID,
		Account: user.Account,
		Version: user.TokenVersion + 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * time.Duration(config.Global.JWT.ExpireDays))),
		}})

	if err != nil {
		// 生成 token 失败，返回 500 错误
		ctx.JSON(500, models.InternalServerError)
		return
	}

	//数据库version++
	_, err = db.SQLEngine.ID(user.ID).Cols("token_version").Update(&db.User{TokenVersion: user.TokenVersion + 1})
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	//redis 缓存 version++
	err = db.SetVersionToCache(user.ID, user.TokenVersion+1)
	if err != nil {
		ctx.JSON(500, models.DatabaseError)
		return
	}

	// 登录成功，返回 token
	ctx.JSON(200, models.LoginResponse{
		Response: models.Success,
		Token:    tokenString,
	})

}
