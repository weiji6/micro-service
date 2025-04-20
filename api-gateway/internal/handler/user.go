package handler

import (
	"context"
	"etcd/api-gateway/internal/service"
	"etcd/api-gateway/pkg/e"
	"etcd/api-gateway/pkg/res"
	"etcd/api-gateway/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegister godoc
// @Summary 用户注册
// @Description 新用户注册接口
// @Tags user
// @Accept json
// @Produce json
// @Param request body service.UserRequest true "用户注册信息"
// @Success 200 {object} res.Response "注册成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 409 {object} res.Response "用户已存在"
// @Router /api/v1/user/register [post]
func UserRegister(ginCtx *gin.Context) {
	var userReq service.UserRequest
	PanicIfUserError(ginCtx.Bind(&userReq))

	// gin.Key 中获取服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserRegister(context.Background(), &userReq)
	PanicIfUserError(err)

	r := res.Response{
		Data:   userResp,
		Status: uint(userResp.Code),
		Msg:    e.GetMsg(uint(userResp.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)
}

// UserLogin godoc
// @Summary 用户登录
// @Description 用户登录接口
// @Tags user
// @Accept json
// @Produce json
// @Param request body service.UserRequest true "用户登录信息"
// @Success 200 {object} res.Response "登录成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 401 {object} res.Response "用户名或密码错误"
// @Router /api/v1/user/login [post]
func UserLogin(ginCtx *gin.Context) {
	var userReq service.UserRequest
	PanicIfUserError(ginCtx.Bind(&userReq))

	// gin.Key 中获取服务实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.UserLogin(context.Background(), &userReq)
	PanicIfUserError(err)

	token, err := util.GenerateToken(uint(userResp.UserDetail.UserID))

	r := res.Response{
		Data: res.TokenData{
			User:  userResp.UserDetail,
			Token: token,
		},
		Status: uint(userResp.Code),
		Msg:    e.GetMsg(uint(userResp.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)
}
