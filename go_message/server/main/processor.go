package main

import (
	"fmt"
	"go_message/common/message"
	"go_message/common/utils"

	// 如果包名和main 里的函数一样，需要重命名包名或函数名
	"go_message/server/process"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 根据客户端发送消息种类不同，决定调用哪个函数来处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	// 测试
	fmt.Println("mes = ", mes)

	switch mes.Type {
	case message.LoginMesType:
		// 处理登录逻辑
		up := &process.UserProcess{
			Conn: (*this).Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		// 处理注册
		up := &process.UserProcess{
			Conn: (*this).Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		// 处理群聊
		smsProcess := &process.SmsProcess{}
		smsProcess.SendGrounpMes(mes)
	default:
		fmt.Println("消息类型不存在， 无法处理...")
	}
	return
}

func (this *Processor) processCtroller() (err error) {

	// 循环读取客户端发送的消息
	for {
		// 封装成函数，返回Message ，Err
		tf := &utils.Transfer{
			Conn: (*this).Conn,
		}
		mes, err := (*tf).ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器也退出...")
				return err
			} else {
				fmt.Println("readPkg err = ", err)
				return err
			}
		}

		fmt.Println("mes = ", mes)

		err = (*this).serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
