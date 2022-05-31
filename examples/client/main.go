package main

import (
	"github.com/zhangdapeng520/zdpgo_nntp"
	"github.com/zhangdapeng520/zdpgo_password"
	"io/ioutil"
	"strings"
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
		panic(err)
	}
	p := zdpgo_password.New(zdpgo_password.PasswordConfig{})
	md5Temp := p.Hash.Md5.EncryptStringNoKey(strings.TrimSpace(string(fileData)))
	client.Log.Debug("上传文件", fileName, md5Temp, "data", string(fileData))

	// 添加文章
	article := &zdpgo_nntp.Article{
		Content: string(fileData),
	}
	err = client.AddArticle(article)
	if err != nil {
		panic(err)
	}

	// 获取文章
	getArticle, err := client.GetArticle(article.Uuid)
	if err != nil {
		panic(err)
	}
	md5Content := p.Hash.Md5.EncryptStringNoKey(strings.TrimSpace(getArticle.Content))
	if md5Temp == md5Content {
		client.Log.Debug("匹配成功", md5Temp, md5Content)
	} else {
		client.Log.Debug("匹配失败", md5Temp, md5Content, "content", getArticle.Content)
	}
}
