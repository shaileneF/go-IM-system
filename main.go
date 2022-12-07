package main

import bean "IM-system/src"

func main() {
	server := bean.NewServer("127.0.0.1", 8888)
	server.Start()
}
