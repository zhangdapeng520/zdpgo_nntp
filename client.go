package zdpgo_nntp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/zhangdapeng520/zdpgo_log"
	"github.com/zhangdapeng520/zdpgo_nntp/cnntp"
	"github.com/zhangdapeng520/zdpgo_password"
	"io/ioutil"
	"time"
)

/*
@Time : 2022/5/31 9:49
@Author : 张大鹏
@File : client.go
@Software: Goland2021.3.1
@Description:
*/

// Client 客户端对象
type Client struct {
	Config *Config
	Log    *zdpgo_log.Log
}

// GetAddress 获取连接地址
func (c *Client) GetAddress() string {
	return fmt.Sprintf("%s:%d",
		c.Config.Client.Host,
		c.Config.Client.Port)
}

// Upload 上传文件
func (c *Client) Upload(username, password string, filedata []byte) string {
	// 连接NNTP服务
	conn, err := cnntp.Dial("tcp", c.GetAddress())
	if err != nil {
		c.Log.Error("获取NNTP连接对象失败", "error", err)
		return ""
	}

	// 权限校验
	if err = conn.Authenticate(username, password); err != nil {
		c.Log.Error("权限校验失败", "error", err)
		return ""
	}

	// 上传POST
	randID := "gXkSUkSclsRlURoYEKFJjtDtsBkEsCsj"
	var newArticle cnntp.Article
	newArticle.Header = map[string][]string{
		"Newsgroups": {"example.all"},                                         // 新闻分组
		"From":       {"<testuser <testuser@example.com>>"},                   // 发送人
		"Subject":    {fmt.Sprintf("Test Subject %s", randID)},                // 标题
		"Date":       {fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))}, // 日期
		"Message-Id": {fmt.Sprintf("<message-%s@example.com>", randID)},       // 消息ID
	}

	// 设置文章的内容
	newArticle.Body = bufio.NewReader(bytes.NewReader(filedata))

	// 上传文章
	err = conn.Post(&newArticle)
	if err != nil {
		fmt.Printf("post new article failed: %v", err)
		// 有时候显示失败，但是还是成功了
		//return ""
	} else {
		fmt.Println("post new article success")
	}

	// 获取分组
	grp := "example.all"
	_, _, _, err = conn.Group(grp)
	if err != nil {
		fmt.Printf("Could not connect to group %s: %v\n", grp, err)
	}

	// 通过ID获取文章
	articleid := fmt.Sprintf("<message-%s@example.com>", randID)
	article, err := conn.Article(articleid)
	if err != nil {
		fmt.Printf("Could not fetch article %s: %v", articleid, err)
	}

	// 读取文章的内容
	body, err := ioutil.ReadAll(article.Body)
	if err != nil {
		fmt.Printf("error reading reader: %v", err)
	}

	// 计算md5
	p := zdpgo_password.New(zdpgo_password.PasswordConfig{})
	// 重新下载的文章末尾会添加换行符，所以去掉
	// 但是某些非文本文件md5还是无法匹配
	md5Temp := p.Hash.Md5.EncryptNoKey(body[0 : len(body)-1])
	fmt.Println(md5Temp)
	return md5Temp

	//没有删除文章的操作
}
