package main

import (
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
	config := zdpgo_nntp.Config{
		Server: zdpgo_nntp.HttpInfo{
			Port: 8887,
		},
	}
	nntp := zdpgo_nntp.NewWithConfig(&config)
	server := nntp.GetServer()
	server.Run()
}
