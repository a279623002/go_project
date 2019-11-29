package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_message/common/message"
	"net"
)

// 将工具函数关联到 Transfer 结构体中
type Transfer struct {
	Conn net.Conn
	// 传输时使用的缓冲数据
	Buf [8192]byte
}

// 解包
func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	fmt.Println("读取客户端发送的数据...")

	// 根据 buf[:4] 读取消息长度
	// conn.Read 在conn 没有被关闭的情况下，才会阻塞
	// 如果客户端关闭了 conn， 就不会阻塞
	_, err = (*this).Conn.Read((*this).Buf[:4])
	if err != nil {
		// fmt.Println("conn.Read err = ", err)
		// 自定义错误
		// err = errors.New("read pkg header error")
		return
	}

	// 根据 buf[:4] 转成一个 uint32 类型（消息长度）
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32((*this).Buf[:4])

	// 根据 pkgLen 读取消息内容
	n, err := (*this).Conn.Read((*this).Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read fail err = ", err)
		return
	}

	// 把 pkgLen 反序列化成 -> message.Message
	// 如果参数 mes 不加&，为空
	err = json.Unmarshal((*this).Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	return
}

// 压包
func (this *Transfer) WritePkg(data []byte) (err error) {
	// 先发送一个长度给对方，防丢包
	// 发送data 数据
	// 先把 data 的长度发给服务器与接下来发送的数据包就行校验（防止丢包）
	// 获取到 data 的长度->转成一个表示长度的 byte 切片
	var pkgLen uint32
	pkgLen = uint32(len(data))

	// 长度为4的原因，一个字节8位，4*8 = 32
	binary.BigEndian.PutUint32((*this).Buf[:4], pkgLen)

	// 发送长度
	// n 为发送的字节数
	n, err := (*this).Conn.Write((*this).Buf[:4])

	if n != 4 || err != nil {
		fmt.Println("conn.write(buf) fail, err = ", err)
		return
	}

	// 发送data 本身
	n, err = (*this).Conn.Write(data)

	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.write(buf) fail, err = ", err)
		return
	}
	return
}
