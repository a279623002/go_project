package process

import (
	"encoding/json"
	"fmt"
	"go_message/common/message"
	"go_message/common/utils"
	"net"
	"os"
)

type UserProcess struct{}

// 注册
func (this *UserProcess) Register(userID int,
	userPwd string, userName string) (err error) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()

	var mes message.Message
	mes.Type = message.RegisterMesType

	var registerMes message.RegisterMes
	registerMes.User.UserID = userID
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(data)

	if err != nil {
		fmt.Println("注册发送信息错误, err = ", err)
		return
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("读取错误 err = ", err)
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登录...")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
}

// 登录
func (this *UserProcess) Login(userID int, userPwd string) (err error) {
	// fmt.Printf("\t 用户ID %d, 密码 %s", userID, userPwd)
	// return nil

	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("ner.Dial err = ", err)
		return
	}
	defer conn.Close()

	// 2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType

	// 3. 创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserID = userID
	loginMes.UserPwd = userPwd

	// 4. 序列化loginMes
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	// 5. 把 data 赋给 mes.Data 字段
	mes.Data = string(data)

	// 6. 序列化mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	// 读取服务器发来的消息
	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg([]byte(data))
	if err != nil {
		fmt.Println("写入错误 err = ", err)
		return
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("读取错误 err = ", err)
		return
	}
	// 将mes 的Data 部分反序列成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// 初始化 CurUser
		CurUser.Conn = conn
		CurUser.User.UserID = userID
		CurUser.User.UserStatus = message.UserOnline

		// 显示当前在线用户列表
		fmt.Println("当前在线用户列表")
		for _, v := range loginResMes.UsersID {
			// 如果不显示自己在线
			if v == userID {
				continue
			}
			fmt.Println("用户ID:", v)

			// 初始化客户端的在线用户
			user := &message.User{
				UserID:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println("\n\n")
		// 在客户端启动一个协程
		// 保持和服务器端的通讯，如果有数据推送给客户端
		// 则接收并显示在客户端
		go serverProcessMes(conn)

		// 循环显示登录成功的菜单，也可以在函数里循环
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}

	return

}
