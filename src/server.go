package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 类似于Java的构造函数
func NewServer(ip string, port int) *Server {
	server := new(Server)
	server.Ip = ip
	server.Port = port
	server.OnlineMap = make(map[string]*User)
	server.Message = make(chan string)
	return server
}

func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
	server.Message <- sendMsg
}

// 监听Message广播消息channel的方法，用协程执行，一旦有消息就发送给全部的在线user
func (server *Server) ListenMessage() {
	for {
		// 不断尝试从这个message管道中读数据，一旦有消息，做后续处理
		// 否则会阻塞在这里，这个管道也是无缓冲的，若读不到数据，就是发送方未发送数据，那么接收方会阻塞
		msg := <-server.Message
		// 将msg发送给全部的在线User
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			user.C <- msg
		}
		server.mapLock.Unlock()
	}
}
func (server *Server) Handler(conn net.Conn) {
	//业务
	fmt.Println("接收到客户端的连接请求", conn)
	//用户上线，将用户加入到onlineMap中
	user := NewUser(conn)
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()
	// 广播当前用户上线消息
	server.BroadCast(user, "已上线")
	select {}
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

	//启动监听Message的goroutine
	go server.ListenMessage()
	for {
		// accept
		// accept之后，代表有客户端在连接此server，说明用户上线
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		// do handler
		go server.Handler(conn)
	}
}
