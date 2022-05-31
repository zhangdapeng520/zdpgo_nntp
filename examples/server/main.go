package main

import (
	"github.com/zhangdapeng520/zdpgo_nntp"
	//"github.com/dustin/go-nntp/server"
	"github.com/zhangdapeng520/zdpgo_nntp/gonntp/server"
	"net"
)

/*
@Time : 2022/5/30 20:34
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

func main() {
	// 获取地址
	a, err := net.ResolveTCPAddr("tcp", ":1119")
	if err != nil {
		panic(err)
	}

	// 创建监听
	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	// 启动服务
	s := nntpserver.NewServer(&zdpgo_nntp.DefaultBackend)

	// 接收客户端信息
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			panic(err)
		}
		go s.Process(c)
	}
}
