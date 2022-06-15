package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_log"
	"github.com/zhangdapeng520/zdpgo_nntp"
)

/*
@Time : 2022/5/30 20:41
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

func main() {
	log := zdpgo_log.NewWithDebug(true, "log.log")
	fileName := "log.log"

	// 获取客户端
	n := zdpgo_nntp.NewWithConfig(&zdpgo_nntp.Config{
		Client: zdpgo_nntp.HttpInfo{
			Host:     "127.0.0.1",
			Port:     8887,
			Username: "zhangdapeng520",
			Password: "zhangdapeng520",
		},
	}, log)
	client := n.GetClient()
	if client.UploadFileAndCheckMd5(fileName) {
		fmt.Println("上传文件成功")
	} else {
		fmt.Println("上传文件失败")
	}
}
