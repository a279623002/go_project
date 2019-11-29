package process

import (
	"encoding/json"
	"fmt"
	"go_message/common/message"
	"go_message/common/utils"
	"net"
)

type SmsProcess struct {
}

// 转发消息
func (this *SmsProcess) SendGrounpMes(mes *message.Message) {

	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("反序列化smsMes失败, err = ", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("反序列化mes失败, err = ", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		// 不转发消息给自己
		if id == smsMes.User.UserID {
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}

}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败, err = ", err)
		return
	}
}
