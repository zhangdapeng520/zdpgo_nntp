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
	n := zdpgo_nntp.NewWithConfig(&zdpgo_nntp.Config{
		Debug: true,
	})
	server := n.GetServer()
	server.Run()
}
