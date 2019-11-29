package process

import (
	"fmt"
	"go_message/client/model"
	"go_message/common/message"
)

var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser // 登录成功后，对其初始化

// 显示在线用户界面
func outputOnlineUser() {
	fmt.Println("当前在线用户列表")
	for id, _ := range onlineUsers {
		fmt.Println("用户ID:\t", id)
	}
}

func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserStatusMes.UserID]
	if !ok {
		// 如果没有则添加
		user = &message.User{
			UserID: notifyUserStatusMes.UserID,
		}
	}

	// 更新状态
	user.UserStatus = notifyUserStatusMes.Status

	onlineUsers[notifyUserStatusMes.UserID] = user

	outputOnlineUser()
}
