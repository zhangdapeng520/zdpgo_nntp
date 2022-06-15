package zdpgo_nntp

import (
	"fmt"
	nntpserver "github.com/zhangdapeng520/zdpgo_nntp/gonntp/server"
	"net"
)

/*
@Time : 2022/5/31 14:45
@Author : 张大鹏
@File : server.go
@Software: Goland2021.3.1
@Description:
*/

// Server NNTP服务
type Server struct {
	Config     *Config
	NntpServer *nntpserver.Server
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
		Log.Error("解析目标地址失败", "error", err)
		return nil, err
	}

	// 创建监听
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		Log.Error("创建监听失败", "error", err)
		return nil, err
	}

	// 返回
	return listener, nil
}

// Handle 处理连接
func (s *Server) Handle(conn net.Conn) {
	s.NntpServer.Process(conn)
}

// Run 运行服务
func (s *Server) Run() error {
	listener, err := s.GetListener()
	if err != nil {
		Log.Panic("创建监听器失败", "error", err)
	}
	defer listener.Close()

	// 接收客户端信息
	var conn net.Conn
	for {
		conn, err = listener.AcceptTCP()
		if err != nil {
			Log.Error("获取客户端连接失败", "error", err)
		}
		go s.Handle(conn)
	}
}
