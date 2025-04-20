package handler

import (
	"context"
	"etcd/task/internal/repository"
	"etcd/task/internal/service"
	"etcd/task/pkg/e"
)

type TaskService struct {
	service.UnimplementedTaskServiceServer
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

// TaskCreate godoc
// @Summary 创建任务
// @Description 创建新的任务
// @Tags task
// @Accept json
// @Produce json
// @Param request body service.TaskRequest true "创建任务请求参数"
// @Success 200 {object} service.CommonResponse "创建成功"
// @Failure 400 {object} service.CommonResponse "请求参数错误"
// @Failure 401 {object} service.CommonResponse "未授权"
// @Router /task/create [post]
func (*TaskService) TaskCreate(ctx context.Context, req *service.TaskRequest) (resp *service.CommonResponse, err error) {
	var task repository.Task
	resp = new(service.CommonResponse)
	resp.Code = e.Success

	err = task.TaskCreate(req)
	if err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}

	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

// TaskShow godoc
// @Summary 查看任务
// @Description 获取任务列表或特定任务详情
// @Tags task
// @Accept json
// @Produce json
// @Param request body service.TaskRequest true "查看任务请求参数"
// @Success 200 {object} service.TaskDetailResponse "成功返回任务信息"
// @Failure 400 {object} service.TaskDetailResponse "请求参数错误"
// @Failure 404 {object} service.TaskDetailResponse "任务不存在"
// @Router /task/show [get]
func (*TaskService) TaskShow(ctx context.Context, req *service.TaskRequest) (resp *service.TaskDetailResponse, err error) {
	var task repository.Task
	resp = new(service.TaskDetailResponse)
	resp.Code = e.Success

	taskList, err := task.TaskShow(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}

	resp.TaskDetail = repository.BuildTasks(taskList)
	return resp, nil
}

// TaskUpdate godoc
// @Summary 更新任务
// @Description 更新现有任务的信息
// @Tags task
// @Accept json
// @Produce json
// @Param request body service.TaskRequest true "更新任务请求参数"
// @Success 200 {object} service.CommonResponse "更新成功"
// @Failure 400 {object} service.CommonResponse "请求参数错误"
// @Failure 404 {object} service.CommonResponse "任务不存在"
// @Router /task/update [put]
func (*TaskService) TaskUpdate(ctx context.Context, req *service.TaskRequest) (resp *service.CommonResponse, err error) {
	var task repository.Task
	resp = new(service.CommonResponse)
	resp.Code = e.Success

	err = task.TaskUpdate(req)
	if err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}

	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

// TaskDelete godoc
// @Summary 删除任务
// @Description 删除指定的任务
// @Tags task
// @Accept json
// @Produce json
// @Param request body service.TaskRequest true "删除任务请求参数"
// @Success 200 {object} service.CommonResponse "删除成功"
// @Failure 400 {object} service.CommonResponse "请求参数错误"
// @Failure 404 {object} service.CommonResponse "任务不存在"
// @Router /task/delete [delete]
func (*TaskService) TaskDelete(ctx context.Context, req *service.TaskRequest) (resp *service.CommonResponse, err error) {
	var task repository.Task
	resp = new(service.CommonResponse)
	resp.Code = e.Success

	err = task.TaskDelete(req)
	if err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}

	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}
