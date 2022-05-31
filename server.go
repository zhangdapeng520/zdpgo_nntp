package zdpgo_nntp

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_log"
	"net"
)

/*
@Time : 2022/5/31 14:45
@Author : 张大鹏
@File : server.go
@Software: Goland2021.3.1
@Description:
*/

type Server struct {
	Config *Config
	Log    *zdpgo_log.Log
}

// GetAddress 获取服务地址
func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d",
		s.Config.Server.Host,
		s.Config.Server.Port)
}
func (s *Server) GetListener() (*net.TCPListener, error) {
	// 获取地址
	addr, err := net.ResolveTCPAddr("tcp", s.GetAddress())
	if err != nil {
		s.Log.Error("解析目标地址失败", "error", err)
		return nil, err
	}

	// 创建监听
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		s.Log.Error("创建监听失败", "error", err)
		return nil, err
	}

	// 返回
	return listener, nil
}
