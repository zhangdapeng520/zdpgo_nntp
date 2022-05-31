package zdpgo_nntp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/zhangdapeng520/zdpgo_log"
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
	Log    *zdpgo_log.Log
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
		c.Log.Error("获取NNTP连接对象失败", "error", err)
		return nil, err
	}

	// 权限校验
	if err = conn.Authenticate(c.Config.Client.Username, c.Config.Client.Password); err != nil {
		c.Log.Error("权限校验失败", "error", err)
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
		c.Log.Error("上传数据失败", "error", err)
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
		c.Log.Error("获取NNTP连接对象失败", "error", err)
		return err
	}

	// 权限校验
	if err = conn.Authenticate(c.Config.Client.Username, c.Config.Client.Password); err != nil {
		c.Log.Error("权限校验失败", "error", err)
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
		c.Log.Error("上传数据失败", "error", err)
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
		c.Log.Error("获取NNTP连接对象失败", "error", err)
		return nil, err
	}

	// 权限校验
	if err = conn.Authenticate(c.Config.Client.Username, c.Config.Client.Password); err != nil {
		c.Log.Error("权限校验失败", "error", err)
		return nil, err
	}

	// 获取分组
	_, _, _, err = conn.Group(c.Config.Group)
	if err != nil {
		c.Log.Error("获取分组失败", "error", err)
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
		c.Log.Error("通过ID获取文章失败", "error", err)
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
		c.Log.Error("读取文章内容失败", "error", err)
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
		c.Log.Error("读取文件失败", "error", err, "filePath", filePath)
		return false
	}
	md5Temp := password.Hash.Md5.EncryptStringNoKey(strings.TrimSpace(string(fileData)))

	// 添加文章
	article := &Article{
		Content: string(fileData),
	}
	err = c.AddArticle(article)
	if err != nil {
		c.Log.Error("上传文件失败", "error", err)
		return false
	}

	// 获取文章
	getArticle, err := c.GetArticle(article.Uuid)
	if err != nil {
		c.Log.Error("获取上传内容失败", "error", err)
		return false
	}
	md5Content := password.Hash.Md5.EncryptStringNoKey(strings.TrimSpace(getArticle.Content))

	// 计算匹配结果并返回
	return md5Temp == md5Content
}
