package main

import (
	"etcd/api-gateway/config"
	"etcd/api-gateway/discovery"
	"etcd/api-gateway/internal/service"
	"etcd/api-gateway/routes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()

	// 服务发现
	etcdAddress := []string{viper.GetString("etcd.address")}
	etcdRegister := discovery.NewResolver(etcdAddress, logrus.New())
	resolver.Register(etcdRegister)

	// 启动 API 网关服务
	go startListen()
	{
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
		s := <-osSignal
		fmt.Println("exit!", s)
	}

	fmt.Println("gateway listen on :4000")
}

func startListen() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	userConn, _ := grpc.Dial(viper.GetString("server.grpcAddress1"), opts...)
	userService := service.NewUserServiceClient(userConn)

	taskConn, _ := grpc.Dial(viper.GetString("server.grpcAddress2"), opts...)
	taskService := service.NewTaskServiceClient(taskConn)

	ginRouter := routes.NewRouter(userService, taskService)
	server := &http.Server{
		Addr:           viper.GetString("server.port"),
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("绑定失败，可能端口被占用", err)
	}
}
