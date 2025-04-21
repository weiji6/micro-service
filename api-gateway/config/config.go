package config

import (
	"github.com/spf13/viper"
)

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
