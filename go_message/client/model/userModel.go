package model

import (
	"go_message/common/message"
	"net"
)

// 在process -> userMgr.go 声明为全局变量
type CurUser struct {
	Conn net.Conn
	message.User
}
