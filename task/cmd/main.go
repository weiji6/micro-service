package main

import (
	"etcd/task/config"
	"etcd/task/discovery"
	"etcd/task/internal/handler"
	"etcd/task/internal/repository"
	"etcd/task/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

// @title User Service API
// @version 1.0
// @description User service for task management system
// @host localhost:8001
// @BasePath /api/v1
func main() {
	config.InitConfig()
	repository.InitDB()

	// etcd 地址
	etcdAddres := []string{viper.GetString("etcd.address")}

	// 服务的注册
	etcdRegister := discovery.NewRegister(etcdAddres, logrus.New())
	grpcAddress := viper.GetString("service.grpcAddress")
	userNode := discovery.Server{
		Name:    viper.GetString("server.domain"),
		Address: grpcAddress,
	}

	server := grpc.NewServer()
	defer server.Stop()

	// 绑定服务
	service.RegisterTaskServiceServer(server, handler.NewTaskService())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(err)
	}
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
