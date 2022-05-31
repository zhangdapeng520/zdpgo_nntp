package zdpgo_nntp

import (
	"bytes"
	"container/ring"
	"github.com/zhangdapeng520/zdpgo_nntp/gonntp"
	nntpserver "github.com/zhangdapeng520/zdpgo_nntp/gonntp/server"
	"io"
	"log"
	"net/textproto"
	"sort"
	"strconv"
	"strings"
)

/*
@Time : 2022/5/30 20:49
@Author : 张大鹏
@File : backend.go
@Software: Goland2021.3.1
@Description:
*/

// 最大文章数量
const maxArticles = 1000

// 文章引用
type articleRef struct {
	msgId string
	num   int64
}

// 分组存储
type groupStorage struct {
	group    *nntp.Group
	articles *ring.Ring
}

// 文章存储
type articleStorage struct {
	headers  textproto.MIMEHeader
	body     string
	refCount int
}

// ServerBackend 服务后台
type ServerBackend struct {
	groups   map[string]*groupStorage   // 分组
	articles map[string]*articleStorage // 文章
	isLogin  bool                       // 是否已登录
}

// DefaultBackend 默认后台
var DefaultBackend = ServerBackend{
	groups:   map[string]*groupStorage{},
	articles: map[string]*articleStorage{},
}

// 初始化
func init() {
	DefaultBackend.groups["zhangdapeng520.all"] = &groupStorage{
		group: &nntp.Group{
			Name:        "zhangdapeng520.all",
			Description: "默认的分组",
			Posting:     nntp.PostingPermitted},
		articles: ring.New(maxArticles),
	}
}

// ListGroups 分组列表
func (tb *ServerBackend) ListGroups(max int) ([]*nntp.Group, error) {
	var rv []*nntp.Group
	for _, g := range tb.groups {
		rv = append(rv, g.group)
	}
	return rv, nil
}

// GetGroup 获取分组
func (tb *ServerBackend) GetGroup(name string) (*nntp.Group, error) {
	var group *nntp.Group

	for _, g := range tb.groups {
		if g.group.Name == name {
			group = g.group
			break
		}
	}

	if group == nil {
		return nil, nntpserver.ErrNoSuchGroup
	}

	return group, nil
}

// mkArticle 创建文章
func mkArticle(a *articleStorage) *nntp.Article {
	return &nntp.Article{
		Header: a.headers,
		Body:   strings.NewReader(a.body),
		Bytes:  len(a.body),
		Lines:  strings.Count(a.body, "\n"),
	}
}

// 在环中查找
func findInRing(in *ring.Ring, f func(r interface{}) bool) *ring.Ring {
	if f(in.Value) {
		return in
	}
	for p := in.Next(); p != in; p = p.Next() {
		if f(p.Value) {
			return p
		}
	}
	return nil
}

// GetArticle 获取文章
func (tb *ServerBackend) GetArticle(group *nntp.Group, id string) (*nntp.Article, error) {

	msgId := id
	var a *articleStorage

	if intId, err := strconv.ParseInt(id, 10, 64); err == nil {
		msgId = ""
		if groupStorage, ok := tb.groups[group.Name]; ok {
			r := findInRing(groupStorage.articles, func(v interface{}) bool {
				if v != nil {
					log.Printf("Looking at %v", v)
				}
				if aref, ok := v.(articleRef); ok && aref.num == intId {
					return true
				}
				return false
			})
			if aref, ok := r.Value.(articleRef); ok {
				msgId = aref.msgId
			}
		}
	}

	a = tb.articles[msgId]
	if a == nil {
		return nil, nntpserver.ErrInvalidMessageID
	}

	return mkArticle(a), nil
}

// 排序的文章列表
type articleList []nntpserver.NumberedArticle

// Len 文章列表长度
func (n articleList) Len() int {
	return len(n)
}

// Less 是否小于
func (n articleList) Less(i, j int) bool {
	return n[i].Num < n[j].Num
}

// Swap 交换
func (n articleList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// GetArticles 获取文章列表
func (tb *ServerBackend) GetArticles(group *nntp.Group,
	from, to int64) ([]nntpserver.NumberedArticle, error) {

	gs, ok := tb.groups[group.Name]
	if !ok {
		Log.Error("获取分组失败", "name", group.Name)
		return nil, nntpserver.ErrNoSuchGroup
	}

	var articles []nntpserver.NumberedArticle
	gs.articles.Do(func(v interface{}) {
		if v != nil {
			if aref, ok := v.(articleRef); ok {
				if aref.num >= from && aref.num <= to {
					a, ok := tb.articles[aref.msgId]
					if ok {
						article := mkArticle(a)
						articles = append(articles,
							nntpserver.NumberedArticle{
								Num:     aref.num,
								Article: article})
					}
				}
			}
		}
	})

	sort.Sort(articleList(articles))
	return articles, nil
}

// AllowPost 是否运行POST提交
func (tb *ServerBackend) AllowPost() bool {
	return true
}

// decr 删除文章
func (tb *ServerBackend) decr(msgId string) {
	if a, ok := tb.articles[msgId]; ok {
		a.refCount--
		if a.refCount == 0 {
			delete(tb.articles, msgId)
		}
	}
}

// Post 新建文章
func (tb *ServerBackend) Post(article *nntp.Article) error {
	// 读取文章内容
	var b []byte
	buf := bytes.NewBuffer(b)
	n, err := io.Copy(buf, article.Body)
	if err != nil {
		Log.Error("读取文章内容失败", "error", err, "length", n)
		return err
	}

	// 查看文章是否已存在
	a := articleStorage{
		headers:  article.Header,
		body:     buf.String(),
		refCount: 0,
	}
	msgId := a.headers.Get("Message-Id")
	if _, ok := tb.articles[msgId]; ok {
		Log.Warning("该文章已存在", "msgId", msgId)
		return nntpserver.ErrPostingFailed
	}

	// 保存到分组
	for _, g := range article.Header["Newsgroups"] {
		if g, ok := tb.groups[g]; ok {
			g.articles = g.articles.Next()
			if g.articles.Value != nil {
				aref := g.articles.Value.(articleRef)
				tb.decr(aref.msgId)
			}
			if g.articles.Value != nil || g.group.Low == 0 {
				g.group.Low++
			}
			g.group.High++
			g.articles.Value = articleRef{
				msgId,
				g.group.High,
			}
			a.refCount++
			g.group.Count = int64(g.articles.Len())
			Log.Debug("保存文章成功", "msgId", msgId, "value", g.articles.Value, "groupName", g.group.Name)
		}
	}

	if a.refCount > 0 {
		tb.articles[msgId] = &a
	} else {
		return nntpserver.ErrPostingFailed
	}

	return nil
}

// Authorized 是否开启权限校验
func (s *ServerBackend) Authorized() bool {
	return s.isLogin
}

// Authenticate 校验用户名和密码
func (s *ServerBackend) Authenticate(user, pass string) (nntpserver.Backend, error) {
	Log.Debug("后台引擎处理权限校验", "user", user, "pass", pass)
	for k, v := range auths {
		if user == k && pass == v.Password {
			s.isLogin = true
			return &DefaultBackend, nil
		}
	}
	return nil, nntpserver.ErrAuthRejected
}
