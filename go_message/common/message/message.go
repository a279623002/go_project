package message

// 定义消息类型常量
const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

type Message struct {
	Type string `json:"type"` // 消息类型，由常量确定
	Data string `json:"data"` // 消息数据
}

// 登录消息
type LoginMes struct {
	UserID   int    `json:"userID"`   // 用户ID
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

// 返回登录消息
type LoginResMes struct {
	Code    int    `json:"code"`  // 返回的状态码 500 表示未注册 200 表示登录成功
	Error   string `json:"error"` // 返回错误信息
	UsersID []int  // 保存在线用户id的切片
}

// 注册消息
type RegisterMes struct {
	User User `json:"user"` // 类型为User结构体
}

// 返回注册消息
type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回的状态码 400 表示用户名已经占有 200 表示注册成功
	Error string `json:"error"` // 返回错误信息
}

// 推送用户上线状态通知消息
type NotifyUserStatusMes struct {
	UserID int `json:"userID"`
	Status int `json:"status"` // 用户状态
}

// 聊天消息
type SmsMes struct {
	User    User   `json:"user"`
	Content string `json:"content"`
}
