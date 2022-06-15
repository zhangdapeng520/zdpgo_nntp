package main

import (
	"github.com/zhangdapeng520/zdpgo_log"
	"github.com/zhangdapeng520/zdpgo_nntp"
)

/*
@Time : 2022/5/30 20:34
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

func main() {
	// 获取监听器
	log := zdpgo_log.NewWithDebug(true, "log.log")
	n := zdpgo_nntp.NewWithConfig(&zdpgo_nntp.Config{
		Server: zdpgo_nntp.HttpInfo{
			Port: 8887,
		},
	}, log)
	server := n.GetServer()
	server.Run()
}
