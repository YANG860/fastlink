package main

import "github.com/gin-gonic/gin"

// getUser 获取用户信息，需验证 token 和权限
func getUser(ctx *gin.Context, account string) {

	var body GetUserRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, GetUserResponse{Response: Response{Success: false, Error: "Invalid request"}})
		return
	}

	// 校验 token
	token, err := ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, GetUserResponse{Response: Response{Success: false, Error: "Invalid token"}})
		return
	}

	// 权限校验
	if token.Account != account {
		ctx.JSON(403, GetUserResponse{Response: Response{Success: false, Error: "Forbidden"}})
		return
	}

	var user User
	// 查询用户
	has, err := engine.ID(token.ID).Get(&user)
	if err != nil {
		ctx.JSON(500, GetUserResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}
	if !has {
		ctx.JSON(404, GetUserResponse{Response: Response{Success: false, Error: "User not found"}})
		return
	}

	// 返回用户信息
	ctx.JSON(200, GetUserResponse{
		Response:        Response{Success: true},
		Account:         user.Account,
		RegisteredAt:    user.RegisteredAt,
		LinkCount:       user.LinkCount,
		LatestCreatedAt: user.LatestCreatedAt,
	})
}
