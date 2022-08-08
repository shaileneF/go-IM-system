package bean

import (
	"fmt"
	"net"
	"strings"
)

// 每个user都会启动一个goroutine，不断监控这个channel，这个channel用于传递消息
// 若监控到这个c里面有消息了，便进行之后的业务处理
// 可以看到在构造方法中，这个c是无缓冲的，说明协程一直监听这个c，若获取不到消息，则会一直阻塞在这里，直到获取到消息
// 对于无缓冲的channel，接收方会阻塞，直到发送方准备好，发送方会阻塞，直到接收方准备好。
type User struct {
	Name   string
	Addr   string
	C      chan string // 这个信道用来接收消息
	conn   net.Conn
	server *Server // 当前用户是属于哪个server的
}

// 构造方法
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := new(User)
	user.Name = userAddr
	user.Addr = userAddr
	user.C = make(chan string)
	user.conn = conn
	user.server = server
	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对应客户端
func (user *User) ListenMessage() {
	for {
		// 一定要有数据，否则会阻塞，这个信道是无缓冲的
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}

// 用户的上线业务
func (user *User) Online() {
	//用户上线，将用户加入到onlineMap中
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	// 广播当前用户上线消息
	user.server.BroadCast(user, "已上线")
}

// 用户的下线业务
func (user *User) Offline() {
	//用户下线，将用户从onlineMap中删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	// 广播当前用户上线消息
	user.server.BroadCast(user, "下线")
}

// 给当前User对应的客户端发消息
//在此方法中的user，是Java中的this,是指的当前对象，也就是调用此方法的对象。
// 也就是谁调用的此方法，就发给谁
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

//用户处理消息的业务
func (user *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前用户都有哪些
		user.server.mapLock.Lock()
		fmt.Println(user.server.OnlineMap)
		for _, u := range user.server.OnlineMap {
			onlineMsg := "[" + u.Addr + "]" + u.Name + ": " + "在线\n"
			// 这里的user是调用sendMsg方法的对象，那么在sendMsg方法中的user，则是Java中的
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		// 判断newName是否存在
		_, isPresent := user.server.OnlineMap[newName]
		if isPresent {
			user.SendMsg("当前用户名被使用\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMsg("您的用户名更新成功:" + user.Name + "\n")
		}
	} else {
		user.server.BroadCast(user, msg)
	}
}
