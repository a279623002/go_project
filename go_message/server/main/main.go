package main

import (
	"fmt"
	"go_message/server/model"
	"net"
	"time"
)

func start(conn net.Conn) {
	// 这里需要延时关闭conn
	defer conn.Close()

	// 调用控制器
	processor := &Processor{
		Conn: conn,
	}

	err := processor.processCtroller()
	if err != nil {
		fmt.Println("客户端和服务器端的协程错误， err = ", err.Error())
		return
	}
}

// 对 UserModel 的初始化
func initUserModel() {
	// 这里的pool (redis.go pool) 本身就是一个全局变量，可在这使用
	// 注意初始化顺序，需要先初始化initPool得到pool
	model.MyUserModel = model.NewUserModel(pool)
}

func init() {
	// 当服务器启动时，初始化连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserModel()
}

func main() {
	// 提示信息
	fmt.Println("服务器在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("ner.Listen err = ", err)
		return
	}

	// 一旦监听成功，就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端来连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err = ", err)
		}

		// 一旦连接成功，启动协程
		go start(conn)
	}
}
