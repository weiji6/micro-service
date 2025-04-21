package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var EtcdClient *clientv3.Client

func InitConfig() {
	viper.SetConfigName("config")   // 配置文件的文件名
	viper.SetConfigType("yaml")     // 配置文件的后缀
	viper.AddConfigPath("./config") // 获取到配置文件的路径
	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置失败：" + err.Error())
	} else {
		println("配置读取成功！使用配置文件:", viper.ConfigFileUsed())
	}
}

func InitEtcdConfig() {
	var err error
	// 读取 etcd 地址
	etcdAddr := viper.GetString("etcd.address")
	if etcdAddr == "" {
		log.Fatal("etcd.address is empty in viper config")
	}

	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}

	loadConfigFromEtcd()
}

func loadConfigFromEtcd() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := "/config/task"
	resp, err := EtcdClient.Get(ctx, key)
	if err != nil {
		log.Fatalf("Cannot get config from etcd: %v", err)
	}

	if len(resp.Kvs) == 0 {
		log.Fatalf("No config found at key: %s", key)
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(resp.Kvs[0].Value))
	if err != nil {
		log.Fatalf("Failed to parse etcd config: %v", err)
	}

	fmt.Println("Loaded config from etcd:", key)
}
