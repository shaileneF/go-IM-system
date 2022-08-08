package bean

import "net"

// 每个user都会启动一个goroutine，不断监控这个channel，这个channel用于传递消息
// 若监控到这个c里面有消息了，便进行之后的业务处理
// 可以看到在构造方法中，这个c是无缓冲的，说明协程一直监听这个c，若获取不到消息，则会一直阻塞在这里，直到获取到消息
// 对于无缓冲的channel，接收方会阻塞，直到发送方准备好，发送方会阻塞，直到接收方准备好。
type User struct {
	Name string
	Addr string
	C    chan string // 这个信道用来接收消息
	conn net.Conn
}

// 构造方法
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := new(User)
	user.Name = userAddr
	user.Addr = userAddr
	user.C = make(chan string)
	user.conn = conn
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
