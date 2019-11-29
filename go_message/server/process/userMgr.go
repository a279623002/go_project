package process

import "fmt"

// 因为UserMgr 实例再服务器端有且只有一个
// 因为在很多的地方，都会使用到，所以定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 初始化
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成对onlineUsers的curd
// 增改一体
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserID] = up
}

// 删除
func (this *UserMgr) DelOnlineUser(userID int) {
	delete(this.onlineUsers, userID)
}

// 返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// 取出用户连接
func (this *UserMgr) GetOnlineUserByID(userID int) (up *UserProcess, err error) {
	// 从mao 取出一个元素，带检测方式
	up, ok := this.onlineUsers[userID]
	if !ok {
		// 说明没找到
		err = fmt.Errorf("用户%d 不在线", userID)
		return
	}
	return
}
