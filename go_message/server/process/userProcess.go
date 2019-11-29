package process

import (
	"encoding/json"
	"fmt"
	"go_message/common/message"
	"go_message/common/utils"
	"go_message/server/model"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	// 表示Conn 是哪个用户
	UserID int
}

// 上线通知
func (this *UserProcess) NotifyOthersOnlineUser(userID int) {
	// 通知所有上线用户
	for id, up := range userMgr.onlineUsers {

		if id == userID {
			// 自己不用通知
			continue
		}

		up.NotifyMeOnline(userID)
	}

}

// 通知上线用户
func (this *UserProcess) NotifyMeOnline(userID int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = userID
	notifyUserStatusMes.Status = message.UserOnline

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("序列化通知数据错误, err = ", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("序列化通知结构体错误, err = ", err)
		return
	}

	tf := &utils.Transfer{
		Conn: (*this).Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送通知结构体错误, err = ", err)
		return
	}

}

// 注册
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	var registerMes message.RegisterMes
	json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("反序列化registerMes失败 err = ", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	var registerResMes message.RegisterResMes

	// 使用model.MyUserModel 验证
	err = model.MyUserModel.Register(&registerMes.User)

	if err != nil {
		if err == message.ERROR_USER_EXITS {
			registerResMes.Code = 505
			registerResMes.Error = message.ERROR_USER_EXITS.Error()
		} else if err == message.ERROR_USER_PWD {
			registerResMes.Code = 506
			registerResMes.Error = "注册发送未知错误..."
		}
	} else {
		registerResMes.Code = 200
	}

	// 序列化registerResMes
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("序列化 registerResMes 失败， err = ", err)
		return
	}

	resMes.Data = string(data)

	// 序列化resMes
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化 resMes 失败， err = ", err)
		return
	}

	// 发送data，封装到 writePkg 函数中
	tf := &utils.Transfer{
		Conn: (*this).Conn,
	}
	err = tf.WritePkg(data)

	return err
}

// 处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1. 先从mes 中取出 mes.Data，并直接反序列化成 LoginMes
	var loginMes message.LoginMes
	json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("反序列化loginMes失败 err = ", err)
		return
	}

	// 声明返回的消息结构体
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes

	// 使用model.MyUserModel 验证
	user, err := model.MyUserModel.Login(loginMes.UserID, loginMes.UserPwd)

	if err != nil {
		if err == message.ERROR_USER_NOTEXITS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == message.ERROR_USER_PWD {
			loginResMes.Code = 300
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = 200
		// 更新用户在线
		// 将登录成功的用户userID 赋给 this
		this.UserID = loginMes.UserID
		userMgr.AddOnlineUser(this)
		// 发起在线通知
		this.NotifyOthersOnlineUser(loginMes.UserID)

		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersID = append(loginResMes.UsersID, id)
		}

		fmt.Println(user, "登录成功")
	}

	// 序列化loginResMes
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("序列化 loginResMes 失败， err = ", err)
		return
	}

	resMes.Data = string(data)

	// 序列化resMes
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化 resMes 失败， err = ", err)
		return
	}

	// 发送data，封装到 writePkg 函数中
	tf := &utils.Transfer{
		Conn: (*this).Conn,
	}
	err = tf.WritePkg(data)

	return err
}
