package main

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := new(Server)
	server.Ip = ip
	server.Port = port
	return server
}
