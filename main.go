package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"

	"chitchat0/common/config"
	"chitchat0/model"
)

var configFile = flag.String("config", "./chitchat0.yaml", "配置文件路径")

// 1. 初始化配置  2. 初始化日志  3. 打开MySql数据库
func init() {
	flag.Parse()

	// 根据指定的文件名*configFile(string类型)初始化配置，信息将保存在config.Conf中，它是config模块中
	// 定义的一个Config类型的实例
	config.InitConfig(*configFile)      // 初始化配置
	initLog()                        // 初始化日志
	err := model.OpenMySql(config.Conf.MySqlUrl, 10, 20, config.Conf.ShowSql, model.Models...) // 连接数据库
	if err != nil {
		log.Error(err)
	}
}

func initLog() {
	// 先按照config.Conf.LogFile的记录打开这个文件，如果没有问题的话就调用logrus.SetOutput将之
	// 设置为日志的输出文件，否则抛出一个错误
	file, err := os.OpenFile(config.Conf.LogFile + "_" + fmt.Sprintf("%d.log", time.Now().Unix()), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)  // 这里file不是一个文件名，而是一个打开了的文件流对象
	} else {
		log.Error(err)
	}
}

func main() {
	initIris()
}
