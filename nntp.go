package zdpgo_nntp

import (
	"github.com/zhangdapeng520/zdpgo_log"
	nntpserver "github.com/zhangdapeng520/zdpgo_nntp/gonntp/server"
	"github.com/zhangdapeng520/zdpgo_password"
)

/*
@Time : 2022/5/30 20:49
@Author : 张大鹏
@File : nntp.go
@Software: Goland2021.3.1
@Description:
*/

var (
	Log      *zdpgo_log.Log
	auths    map[string]Auth
	password = zdpgo_password.New(Log)
)

type Nntp struct {
	Config *Config
	Log    *zdpgo_log.Log
}

func New(log *zdpgo_log.Log) *Nntp {
	return NewWithConfig(&Config{}, log)
}

func NewWithConfig(config *Config, log *zdpgo_log.Log) *Nntp {
	n := &Nntp{}

	// 日志
	Log = log
	nntpserver.Log = log

	// 配置
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 35333
	}
	if config.Client.Host == "" {
		config.Client.Host = "127.0.0.1"
	}
	if config.Client.Port == 0 {
		config.Client.Port = 35333
	}
	if config.Client.Username == "" {
		config.Client.Username = "zhangdapeng520"
	}
	if config.Client.Password == "" {
		config.Client.Password = "zhangdapeng520"
	}
	if config.Auths == nil || len(config.Auths) == 0 {
		config.Auths = map[string]Auth{
			"zhangdapeng520": {"zhangdapeng520", "zhangdapeng520"},
		}
	}
	if config.Group == "" {
		config.Group = "zhangdapeng520.all"
	}
	if config.From == "" {
		config.From = "<zhangdapeng520 <zhangdapeng520@zdpgo.com>>"
	}
	n.Config = config

	// 权限数据
	auths = n.Config.Auths

	// 返回
	return n
}

// GetClient 获取客户端
func (n *Nntp) GetClient() *Client {
	return &Client{
		Config: n.Config,
	}
}

// GetServer 获取服务对象
func (n *Nntp) GetServer() *Server {
	return &Server{
		Config:     n.Config,
		NntpServer: nntpserver.NewServer(&DefaultBackend),
	}
}
