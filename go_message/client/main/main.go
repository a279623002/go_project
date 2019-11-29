package main

import (
	"fmt"
	"go_message/client/process"
	"os"
)

var (
	userID   int
	userPwd  string
	userName string
)

func main() {
	var key int

	for {
		fmt.Println("------------海量用户通讯系统--------------")
		fmt.Println("\t 1. 登录聊天室")
		fmt.Println("\t 2. 注册用户")
		fmt.Println("\t 3. 退出系统")

		fmt.Println("\t 登录聊天室")

		fmt.Println("\t 请选择（1-3）")
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("\t 登录聊天室")
			fmt.Println("\t 请输入用户ID")
			fmt.Scanf("%d\n", &userID)
			fmt.Println("\t 请输入用户密码")
			fmt.Scanf("%s\n", &userPwd)

			// 创建UserProcess的实例,调用登录方法
			up := &process.UserProcess{}

			err := up.Login(userID, userPwd)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			fmt.Println("注册用户")
			fmt.Println("\t 请输入用户ID")
			fmt.Scanf("%d\n", &userID)
			fmt.Println("\t 请输入用户密码")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("\t 请输入用户名")
			fmt.Scanf("%s\n", &userName)

			// 创建UserProcess的实例,调用注册方法
			up := &process.UserProcess{}

			up.Register(userID, userPwd, userName)

		case 3:
			fmt.Println("退出系统")
			// 结束程序
			os.Exit(0)
		default:
			fmt.Println("输入错误，请重试")
		}
	}
}
