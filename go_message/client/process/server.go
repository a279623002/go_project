package process

import (
	"encoding/json"
	"fmt"
	"go_message/common/message"
	"go_message/common/utils"
	"net"
	"os"
)

// 显示登录成功后的界面
func ShowMenu() {
	fmt.Println("---------登录成功--------")
	fmt.Println("---------1. 显示在线用户列表--------")
	fmt.Println("---------2. 发送消息--------")
	fmt.Println("---------3. 信息列表--------")
	fmt.Println("---------4. 退出系统--------")
	fmt.Println("请选择（1-4）：")
	var key int
	var content string
	SmsProcess := &SmsProcess{}

	fmt.Scanf("%d\n", &key)

	switch key {
	case 1:
		outputOnlineUser()
	case 2:
		fmt.Println("你想对大家说点什么 :)")
		fmt.Scanf("%s\n", &content)
		SmsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		// 结束程序
		os.Exit(0)
	default:
		fmt.Println("输入错误")
	}
}

// 和服务器保持通讯，处理推送的消息
func serverProcessMes(conn net.Conn) {
	// 创建transfer 实例，不停读取服务器发送的消息
	// 前提，连接畅通未被关闭
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("服务器错误, tf.ReadPkg err = ", err)
			return
		}
		// 如果读取到消息，处理逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 上线通知
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			// 群发消息
			outputGrounpMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")
		}
	}
}
