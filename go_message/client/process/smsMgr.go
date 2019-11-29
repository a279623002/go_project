package process

import (
	"encoding/json"
	"fmt"
	"go_message/common/message"
)

func outputGrounpMes(mes *message.Message) {
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("反序列化mes失败, err = ", err.Error())
		return
	}

	info := fmt.Sprintf("用户ID：\t%d 对大家说：\t%s", smsMes.User.UserID, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
