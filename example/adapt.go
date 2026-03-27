package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyahm/xmux" // 改成你的包名
)

// ====================== 业务流程：用户登录（带分支）======================
type UserLoginFlow struct {
	xmux.BaseFlow // 继承引擎
	// 业务字段
	Username string
	Password string
	IsAdmin  bool // 分支条件
}

// 构造函数
func NewUserLoginFlow() xmux.Flow {
	return &UserLoginFlow{}
}

// ====================== ✅ 终极纯链式写法 ======================
func (f *UserLoginFlow) Run() {
	f.
		Then(f.ParseParam).    // 1. 解析参数
		Then(f.CheckPassword). // 2. 校验密码

		// 🔹 分支：管理员 / 普通用户
		IfElse(f.IsAdmin, f.CheckAdminAuth, f.CheckUserAuth).
		Then(f.CreateToken).    // 3. 生成令牌
		Then(f.ResponseSuccess) // 4. 返回成功

	// 统一错误处理
	if f.Err() != nil {
		f.ResponseFail(f.Err())
	}
}

// ====================== 业务步骤 ======================
// 解析参数
func (f *UserLoginFlow) ParseParam() error {
	f.Username = f.R.URL.Query().Get("username")
	f.Password = f.R.URL.Query().Get("password")

	if f.Username == "" || f.Password == "" {
		return errors.New("用户名或密码不能为空")
	}

	// 模拟管理员
	if f.Username == "admin" {
		f.IsAdmin = true
	}
	return nil
}

// 校验密码
func (f *UserLoginFlow) CheckPassword() error {
	if f.Password != "123456" {
		return errors.New("密码错误")
	}
	return nil
}

// 管理员校验
func (f *UserLoginFlow) CheckAdminAuth() error {
	fmt.Println("执行【管理员】权限校验")
	return nil
}

// 普通用户校验
func (f *UserLoginFlow) CheckUserAuth() error {
	fmt.Println("执行【普通用户】权限校验")
	return nil
}

// 生成Token
func (f *UserLoginFlow) CreateToken() error {
	f.Ins.Set("token", "token_"+f.Username)
	return nil
}

// 响应成功
func (f *UserLoginFlow) ResponseSuccess() error {
	_ = json.NewEncoder(f.W).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "登录成功",
		"data": f.Ins.Data,
	})
	return nil
}

// 响应失败
func (f *UserLoginFlow) ResponseFail(err error) {
	_ = json.NewEncoder(f.W).Encode(map[string]interface{}{
		"code": -1,
		"msg":  err.Error(),
	})
}

// ====================== 路由启动 ======================
func main() {
	router := xmux.NewRouter()
	router.Get("/login", xmux.Adapt(NewUserLoginFlow))
	router.Run()
}
