package handler

import (
	"context"
	"etcd/user/internal/repository"
	"etcd/user/internal/service"
	"etcd/user/pkg/e"
)

type UserService struct {
	service.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

// UserLogin godoc
// @Summary 用户登录
// @Description 处理用户登录请求，验证用户凭据并返回用户详细信息
// @Tags user
// @Accept json
// @Produce json
// @Param request body service.UserRequest true "用户登录请求参数"
// @Success 200 {object} service.UserDetailResponse "成功返回用户详细信息"
// @Failure 400 {object} service.UserDetailResponse "请求参数错误"
// @Failure 401 {object} service.UserDetailResponse "用户验证失败"
// @Router /user/login [post]
func (*UserService) UserLogin(ctx context.Context, req *service.UserRequest) (resp *service.UserDetailResponse, err error) {
	var user repository.User
	resp = new(service.UserDetailResponse)
	resp.Code = e.Success
	err = user.ShowUserInfo(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}

	resp.UserDetail = repository.BuildUser(user)
	return resp, nil
}

// UserRegister godoc
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags user
// @Accept json
// @Produce json
// @Param request body service.UserRequest true "用户注册请求参数"
// @Success 200 {object} service.UserDetailResponse "注册成功返回用户信息"
// @Failure 400 {object} service.UserDetailResponse "请求参数错误"
// @Failure 409 {object} service.UserDetailResponse "用户已存在"
// @Router /user/register [post]
func (*UserService) UserRegister(ctx context.Context, req *service.UserRequest) (resp *service.UserDetailResponse, err error) {
	var user repository.User
	resp = new(service.UserDetailResponse)
	resp.Code = e.Success
	user, err = user.UserCreate(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}

	resp.UserDetail = repository.BuildUser(user)
	return resp, nil
}
