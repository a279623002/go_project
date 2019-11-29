package model

import (
	"encoding/json"
	"fmt"
	"go_message/common/message"

	"github.com/garyburd/redigo/redis"
)

// 在服务器启动后，初始化一个 userModel 实例
// 把它做成全局的变量，在需要和redis 操作时，直接使用即可
var (
	MyUserModel *UserModel
)

// 定义一个结构体，对User 结构体的CURD
type UserModel struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个 UserModel 实例
func NewUserModel(pool *redis.Pool) (userModel *UserModel) {
	userModel = &UserModel{
		pool: pool,
	}
	return
}

// 根据用户id 返回User实例、err
func (this *UserModel) getUserByID(conn redis.Conn, id int) (user message.User, err error) {
	// 通过 id 去 redis 查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		// 调用自定义错误
		if err == redis.ErrNil {
			// 表示在 users 哈希中，没有找到对应id
			err = message.ERROR_USER_NOTEXITS
		}
		return
	}

	// 把结构反序列成 User 实例
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("json 反序列化错误 err = ", err)
		return
	}
	return
}

// 完成登录的校验 Login
// 如果用户的id 和pwd 都正确，则返回一个user 实例
// 如果用户的id 和pwd 不正确，返回对应的错误信息
func (this *UserModel) Login(userID int, userPwd string) (user message.User, err error) {
	// 先从 UserModel 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserByID(conn, userID)
	if err != nil {
		return
	}

	// 用户存在，校验密码
	if user.UserPwd != userPwd {
		err = message.ERROR_USER_PWD
		return
	}

	return
}

// 注册
func (this *UserModel) Register(user *message.User) (err error) {
	// 先从 UserModel 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserByID(conn, user.UserID)
	if err == nil {
		// 用户已存在，注册错误
		err = message.ERROR_USER_EXITS
		return
	}

	// 序列化user
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	// 入库
	_, err = conn.Do("HSet", "users", user.UserID, string(data))
	if err != nil {
		fmt.Println("注册错误 err = ", err)
		return
	}
	return
}
