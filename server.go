package zdpgo_nntp

import (
	"bytes"
	"container/ring"
	"github.com/zhangdapeng520/zdpgo_nntp/gonntp"
	//"github.com/dustin/go-nntp"
	//nntpserver "github.com/dustin/go-nntp/server"
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
@File : server.go
@Software: Goland2021.3.1
@Description:
*/

// 最大文章数量
const maxArticles = 1000

// 文章引用
type articleRef struct {
	msgid string
	num   int64
}

// 分组存储
type groupStorage struct {
	group *nntp.Group
	// article refs
	articles *ring.Ring
}

// 文章存储
type articleStorage struct {
	headers  textproto.MIMEHeader
	body     string
	refcount int
}

// ServerBackend 服务后台
type ServerBackend struct {
	// group name -> group storage
	groups map[string]*groupStorage
	// message ID -> article
	articles map[string]*articleStorage
}

// DefaultBackend 默认后台
var DefaultBackend = ServerBackend{
	groups:   map[string]*groupStorage{},
	articles: map[string]*articleStorage{},
}

// 初始化
func init() {
	DefaultBackend.groups["example.all"] = &groupStorage{
		group: &nntp.Group{
			Name:        "example.all",
			Description: "A test.",
			Posting:     nntp.PostingPermitted},
		articles: ring.New(maxArticles),
	}
}

// ListGroups 分组列表
func (tb *ServerBackend) ListGroups(max int) ([]*nntp.Group, error) {
	rv := []*nntp.Group{}
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

	msgID := id
	var a *articleStorage

	if intid, err := strconv.ParseInt(id, 10, 64); err == nil {
		msgID = ""
		// by int ID.  Gotta go find it.
		if groupStorage, ok := tb.groups[group.Name]; ok {
			r := findInRing(groupStorage.articles, func(v interface{}) bool {
				if v != nil {
					log.Printf("Looking at %v", v)
				}
				if aref, ok := v.(articleRef); ok && aref.num == intid {
					return true
				}
				return false
			})
			if aref, ok := r.Value.(articleRef); ok {
				msgID = aref.msgid
			}
		}
	}

	a = tb.articles[msgID]
	if a == nil {
		return nil, nntpserver.ErrInvalidMessageID
	}

	return mkArticle(a), nil
}

// 排序的文章列表
type nalist []nntpserver.NumberedArticle

// Len 文章列表长度
func (n nalist) Len() int {
	return len(n)
}

// Less 是否小于
func (n nalist) Less(i, j int) bool {
	return n[i].Num < n[j].Num
}

// Swap 交换
func (n nalist) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// GetArticles 获取文章列表
func (tb *ServerBackend) GetArticles(group *nntp.Group,
	from, to int64) ([]nntpserver.NumberedArticle, error) {

	gs, ok := tb.groups[group.Name]
	if !ok {
		return nil, nntpserver.ErrNoSuchGroup
	}

	log.Printf("Getting articles from %d to %d", from, to)

	rv := []nntpserver.NumberedArticle{}
	gs.articles.Do(func(v interface{}) {
		if v != nil {
			if aref, ok := v.(articleRef); ok {
				if aref.num >= from && aref.num <= to {
					a, ok := tb.articles[aref.msgid]
					if ok {
						article := mkArticle(a)
						rv = append(rv,
							nntpserver.NumberedArticle{
								Num:     aref.num,
								Article: article})
					}
				}
			}
		}
	})

	sort.Sort(nalist(rv))

	return rv, nil
}

// AllowPost 是否运行POST提交
func (tb *ServerBackend) AllowPost() bool {
	return true
}

// decr 删除文章
func (tb *ServerBackend) decr(msgid string) {
	if a, ok := tb.articles[msgid]; ok {
		a.refcount--
		if a.refcount == 0 {
			log.Printf("Getting rid of %v", msgid)
			delete(tb.articles, msgid)
		}
	}
}

// Post 新建文章
func (tb *ServerBackend) Post(article *nntp.Article) error {
	log.Printf("Got headers: %#v", article.Header)
	b := []byte{}
	buf := bytes.NewBuffer(b)
	n, err := io.Copy(buf, article.Body)
	if err != nil {
		return err
	}
	log.Printf("Read %d bytes of body", n)

	a := articleStorage{
		headers:  article.Header,
		body:     buf.String(),
		refcount: 0,
	}

	msgID := a.headers.Get("Message-Id")

	if _, ok := tb.articles[msgID]; ok {
		return nntpserver.ErrPostingFailed
	}

	for _, g := range article.Header["Newsgroups"] {
		if g, ok := tb.groups[g]; ok {
			g.articles = g.articles.Next()
			if g.articles.Value != nil {
				aref := g.articles.Value.(articleRef)
				tb.decr(aref.msgid)
			}
			if g.articles.Value != nil || g.group.Low == 0 {
				g.group.Low++
			}
			g.group.High++
			g.articles.Value = articleRef{
				msgID,
				g.group.High,
			}
			log.Printf("Placed %v", g.articles.Value)
			a.refcount++
			g.group.Count = int64(g.articles.Len())

			log.Printf("Stored %v in %v", msgID, g.group.Name)
		}
	}

	if a.refcount > 0 {
		tb.articles[msgID] = &a
	} else {
		return nntpserver.ErrPostingFailed
	}

	return nil
}

// Authorized 是否开启权限校验
func (tb *ServerBackend) Authorized() bool {
	return true
}

// Authenticate 校验用户名和密码
func (tb *ServerBackend) Authenticate(user, pass string) (nntpserver.Backend, error) {
	return nil, nntpserver.ErrAuthRejected
}

// 出错了
func maybefatal(err error, f string, a ...interface{}) {
	if err != nil {
		log.Fatalf(f, a...)
	}
}
