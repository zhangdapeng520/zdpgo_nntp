package main

import (
	"github.com/zhangdapeng520/zdpgo_nntp"
	//"github.com/dustin/go-nntp/server"
	"github.com/zhangdapeng520/zdpgo_nntp/gonntp/server"
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
	listener, err := server.GetListener()
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// 启动服务
	s := nntpserver.NewServer(&zdpgo_nntp.DefaultBackend)

	// 接收客户端信息
	for {
		c, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		go s.Process(c)
	}
}
