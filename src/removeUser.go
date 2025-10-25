package main

import "github.com/gin-gonic/gin"

// removeUser 注销用户（逻辑删除），需验证 token 和权限
func removeUser(ctx *gin.Context, account string) {

	var body RemoveUserRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, RemoveUserResponse{Response: Response{Success: false, Error: "Invalid request"}})
		return
	}

	// 校验 token
	userToken, err := ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, RemoveUserResponse{Response: Response{Success: false, Error: "Invalid token"}})
		return
	}

	// 权限校验
	if userToken.Account != account {
		ctx.JSON(403, RemoveUserResponse{Response: Response{Success: false, Error: "Forbidden"}})
		return
	}

	// 逻辑删除用户
	_, err = engine.ID(userToken.ID).Cols("valid").Update(&User{Valid: false})
	if err != nil {
		ctx.JSON(500, RemoveUserResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}
	ctx.JSON(200, RemoveUserResponse{Response: Response{Success: true}})
}
