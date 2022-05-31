package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_nntp"
	"github.com/zhangdapeng520/zdpgo_password"
	"github.com/zhangdapeng520/zdpgo_uuid"
	"io/ioutil"
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

	// 计算md5
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("read fail", err)
	}
	p := zdpgo_password.New(zdpgo_password.PasswordConfig{})
	md5Temp := p.Hash.Md5.EncryptNoKey(fileData)
	fmt.Printf("文件 %s md5=%s\n", fileName, md5Temp)

	response, err := client.PostBytes(zdpgo_uuid.UUID(), fileData)
	if err != nil {
		panic(err)
	}
	if md5Temp == response.Text {
		fmt.Printf("上传文件 %s 成功，上传方式 NNTP", fileName)
	} else {
		fmt.Printf("上传文件 %s 失败，上传方式 NNTP,MD5不匹配", fileName)
	}
}
