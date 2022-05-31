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
	fileName := "test.txt"

	// 获取客户端
	n := zdpgo_nntp.NewWithConfig(&zdpgo_nntp.Config{
		Debug: true,
	})
	client := n.GetClient()
	if client.UploadFileAndCheckMd5(fileName) {
		fmt.Println("上传文件成功")
	} else {
		fmt.Println("上传文件失败")
	}
}
