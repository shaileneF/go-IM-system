package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 类似于Java的构造函数
func NewServer(ip string, port int) *Server {
	server := new(Server)
	server.Ip = ip
	server.Port = port
	return server
}

func (server *Server) Handler(conn net.Conn) {
	//业务
	fmt.Println("接收到客户端的连接请求")
}

//启动服务器的函数
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		// do handler
		go server.Handler(conn)
	}
}
