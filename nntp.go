package zdpgo_nntp

import "github.com/zhangdapeng520/zdpgo_log"

/*
@Time : 2022/5/30 20:49
@Author : 张大鹏
@File : nntp.go
@Software: Goland2021.3.1
@Description:
*/

type Nntp struct {
	Config *Config
	Log    *zdpgo_log.Log
}

func New() *Nntp {
	return NewWithConfig(&Config{})
}

func NewWithConfig(config *Config) *Nntp {
	n := &Nntp{}

	// 日志
	if config.LogFilePath == "" {
		config.LogFilePath = "logs/zdpgo/zdpgo_nntp.log"
	}
	n.Log = zdpgo_log.NewWithDebug(config.Debug, config.LogFilePath)

	// 配置
	n.Config = config

	// 返回
	return n
}

func (n *Nntp) GetClient() *Client {
	return &Client{
		Config: n.Config,
		Log:    n.Log,
	}
}
