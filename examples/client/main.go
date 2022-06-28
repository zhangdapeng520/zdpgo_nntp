package main

import (
	"fmt"
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
	// 获取客户端
	n := zdpgo_nntp.NewWithConfig(&zdpgo_nntp.Config{
		Client: zdpgo_nntp.HttpInfo{
			Host:     "127.0.0.1",
			Port:     8887,
			Username: "zhangdapeng520",
			Password: "zhangdapeng520",
		},
	})
	client := n.GetClient()
	if client.UploadFileAndCheckMd5("README.md") {
		fmt.Println("上传文件成功")
	} else {
		fmt.Println("上传文件失败")
	}
}
