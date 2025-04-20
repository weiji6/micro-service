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

// ListTask 获取所有事件
// @Summary 获取所有事件
// @Description 用户获取所有事件
// @Tags 事件
// @Accept json
// @Produce json
// @Success 200 {object} res.Response "获取成功"
// @Failure 400 {object} res.Response "解析失败"
// @Failure 401 {object} res.Response "用户未授权"
// @Failure 500 {object} res.Response "服务器错误"
// @Router /v1/api/task [get]
func ListTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	PanicIfTaskError(ginCtx.Bind(&tReq))

	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskService := ginCtx.Keys["task"].(service.TaskServiceClient)
	taskResp, err := taskService.TaskShow(context.Background(), &tReq)
	PanicIfUserError(err)

	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)
}

// CreateTask 创建新事件
// @Summary 创建新事件
// @Description 用户创建新事件
// @Tags 事件
// @Accept json
// @Produce json
// @Param request body service.TaskRequest true "事件数据"
// @Success 200 {object} res.Response "获取成功"
// @Failure 400 {object} res.Response "解析失败"
// @Failure 401 {object} res.Response "用户未授权"
// @Failure 500 {object} res.Response "服务器错误"
// @Router /v1/api/task [post]
func CreateTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	PanicIfTaskError(ginCtx.Bind(&tReq))

	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskService := ginCtx.Keys["task"].(service.TaskServiceClient)
	taskResp, err := taskService.TaskCreate(context.Background(), &tReq)
	PanicIfUserError(err)

	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)
}

// UpdateTask 更新事件
// @Summary 更新事件
// @Description 用户更新事件
// @Tags 事件
// @Accept json
// @Produce json
// @Param request body service.TaskRequest true "事件数据"
// @Success 200 {object} res.Response "获取成功"
// @Failure 400 {object} res.Response "解析失败"
// @Failure 401 {object} res.Response "用户未授权"
// @Failure 500 {object} res.Response "服务器错误"
// @Router /v1/api/task [put]
func UpdateTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	PanicIfTaskError(ginCtx.Bind(&tReq))

	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskService := ginCtx.Keys["task"].(service.TaskServiceClient)
	taskResp, err := taskService.TaskUpdate(context.Background(), &tReq)
	PanicIfUserError(err)

	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)
}

// DeleteTask 删除事件
// @Summary 删除事件
// @Description 用户删除事件
// @Tags 事件
// @Accept json
// @Produce json
// @Param request body service.TaskRequest true "事件数据"
// @Success 200 {object} res.Response "获取成功"
// @Failure 400 {object} res.Response "解析失败"
// @Failure 401 {object} res.Response "用户未授权"
// @Failure 500 {object} res.Response "服务器错误"
// @Router /v1/api/task [delete]
func DeleteTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	PanicIfTaskError(ginCtx.Bind(&tReq))

	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskService := ginCtx.Keys["task"].(service.TaskServiceClient)
	taskResp, err := taskService.TaskDelete(context.Background(), &tReq)
	PanicIfUserError(err)

	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)
}
