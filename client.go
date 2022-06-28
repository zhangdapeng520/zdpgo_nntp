package zdpgo_nntp

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/zhangdapeng520/zdpgo_nntp/cnntp"
	"github.com/zhangdapeng520/zdpgo_uuid"
	"io/ioutil"
	"strings"
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
}

// GetAddress 获取连接地址
func (c *Client) GetAddress() string {
	return fmt.Sprintf("%s:%d",
		c.Config.Client.Host,
		c.Config.Client.Port)
}

// PostBytes 上传字节数组
func (c *Client) PostBytes(title string, data []byte) (*Response, error) {
	response := &Response{}

	// 连接NNTP服务
	conn, err := cnntp.Dial("tcp", c.GetAddress())
	if err != nil {
		return nil, err
	}

	// 权限校验
	if err = conn.Authenticate(c.Config.Client.Username, c.Config.Client.Password); err != nil {
		return nil, err
	}

	// 上传POST
	uuid := zdpgo_uuid.UUID()
	response.Uuid = uuid
	if title == "" {
		title = uuid
	}
	var newArticle cnntp.Article
	newArticle.Header = map[string][]string{
		"Newsgroups": {c.Config.Group},                                        // 新闻分组
		"From":       {c.Config.From},                                         // 发送人
		"Subject":    {title},                                                 // 标题
		"Date":       {fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))}, // 日期
		"Message-Id": {fmt.Sprintf("<message-%s@zdpgo.com>", uuid)},           // 消息ID
	}

	// 设置文章的内容
	newArticle.Body = bufio.NewReader(bytes.NewReader(data))

	// 上传文章
	err = conn.Post(&newArticle)
	if err != nil {
		return nil, err
	}

	response.StatusCode = 200
	return response, nil
}

// AddArticle 添加文章
func (c *Client) AddArticle(article *Article) error {
	// 连接NNTP服务
	conn, err := cnntp.Dial("tcp", c.GetAddress())
	if err != nil {
		return err
	}

	// 权限校验
	if err = conn.Authenticate(c.Config.Client.Username, c.Config.Client.Password); err != nil {
		return err
	}

	// 上传POST
	if article.Uuid == "" {
		article.Uuid = zdpgo_uuid.UUID()
	}
	if article.Title == "" {
		article.Title = article.Uuid
	}
	if article.Group == "" {
		article.Group = c.Config.Group
	}
	if article.Author == "" {
		article.Author = c.Config.From
	}
	article.Date = int(time.Now().Unix())
	article.DateTime = time.Now()
	article.DateString = time.Now().Format("2006-01-02 15:04:05")
	article.MessageId = fmt.Sprintf("<message-%s@zdpgo.com>", article.Uuid)

	var newArticle cnntp.Article
	newArticle.Header = map[string][]string{
		"Newsgroups": {article.Group},      // 新闻分组
		"From":       {article.Author},     // 发送人
		"Subject":    {article.Title},      // 标题
		"Date":       {article.DateString}, // 日期
		"Message-Id": {article.MessageId},  // 消息ID
	}

	// 设置文章的内容
	newArticle.Body = bufio.NewReader(bytes.NewReader([]byte(article.Content)))

	// 上传文章
	err = conn.Post(&newArticle)
	if err != nil {
		return err
	}

	// 返回
	return nil
}

// GetArticle 根据UUID获取文章
func (c *Client) GetArticle(uuid string) (*Article, error) {
	// 连接NNTP服务
	conn, err := cnntp.Dial("tcp", c.GetAddress())
	if err != nil {
		return nil, err
	}

	// 权限校验
	if err = conn.Authenticate(c.Config.Client.Username, c.Config.Client.Password); err != nil {
		return nil, err
	}

	// 获取分组
	_, _, _, err = conn.Group(c.Config.Group)
	if err != nil {
		return nil, err
	}

	// 创建文章
	article := &Article{
		Uuid: uuid,
	}

	// 通过ID获取文章
	articleId := fmt.Sprintf("<message-%s@zdpgo.com>", uuid)
	nntpArticle, err := conn.Article(articleId)
	if err != nil {
		return nil, err
	}

	// 解析文章
	article.Group = nntpArticle.Header["Newsgroups"][0]
	article.Author = nntpArticle.Header["From"][0]
	article.Title = nntpArticle.Header["Subject"][0]
	article.DateString = nntpArticle.Header["Date"][0]
	article.MessageId = nntpArticle.Header["Message-Id"][0]

	// 解析日期
	timeLayout := "2006-01-02 15:04:05"                                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                                            //重要：获取时区
	article.DateTime, _ = time.ParseInLocation(timeLayout, article.DateString, loc) //使用模板在对应时区转化为time.time类型
	article.Date = int(article.DateTime.Unix())

	// 读取文章的内容
	body, err := ioutil.ReadAll(nntpArticle.Body)
	if err != nil {
		return nil, err
	}
	article.Content = string(body)

	// 返回
	return article, nil
}

// UploadFileAndCheckMd5 上传文件并检查MD5
func (c *Client) UploadFileAndCheckMd5(filePath string) bool {
	// 读取文件
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}
	// 移除空格
	tempDta := strings.TrimSpace(string(fileData))
	tempDta = strings.Replace(tempDta, "\r\n", "", -1)
	tempDta = strings.Replace(tempDta, "\n", "", -1)
	md5Temp := c.Md5(tempDta)

	// 添加文章
	article := &Article{
		Content: string(fileData),
	}
	err = c.AddArticle(article)
	if err != nil {
		return false
	}

	// 获取文章
	getArticle, err := c.GetArticle(article.Uuid)
	if err != nil {
		return false
	}
	contentData := strings.TrimSpace(getArticle.Content)
	contentData = strings.Replace(contentData, "\r\n", "", -1)
	contentData = strings.Replace(contentData, "\n", "", -1)
	md5Content := c.Md5(contentData)

	// 计算匹配结果并返回
	return md5Temp == md5Content
}

// Md5 校验MD5
func (c *Client) Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
