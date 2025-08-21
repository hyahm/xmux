package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	Id          int64  `xorm:"not null pk autoincr comment('用户ID') INT(11)"`
	Username    string `xorm:"not null unique comment('用户账号') VARCHAR(50)" json:"username"`
	Nickname    string `xorm:"not null unique comment('用户昵称') VARCHAR(32)"`
	Avatar      string `xorm:"default null comment('用户头像url') VARCHAR(500)"`
	Background  string `xorm:"default null comment('主页背景图url') VARCHAR(500)"`
	Gender      int8   `xorm:"not null default 2 comment('性别 0女 1男 2未知') TINYINT(4)"`
	Description string `xorm:"default null comment('个性签名') VARCHAR(100)"`
	Vip         int8   `xorm:"not null default 0 comment('会员类型 0普通用户 1月度大会员 2季度大会员 3年度大会员') TINYINT(4)"`
	Role        int8   `xorm:"not null default 0 comment('角色类型 0普通用户 1管理员 2超级管理员') TINYINT(4)"`
	Auth        int8   `xorm:"not null default 0 comment('官方认证 0普通用户 1个人认证 2机构认证') TINYINT(4)"`
	AuthMsg     string `xorm:"default null comment('认证说明') VARCHAR(30)"`
}

var jwtKey = []byte("mysecretsigningkey") // 生产环境中应从环境变量或配置文件中读取

type MyCustomClaims struct {
	UserInfo UserInfo
	jwt.RegisteredClaims
}

func GenerateToken(info UserInfo) (string, error) {

	claims := MyCustomClaims{
		UserInfo: info,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

var ErrParseToken = errors.New("error parsing token")

func VerifyToken(tokenString string) (UserInfo, error) {
	var claims MyCustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return UserInfo{}, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.UserInfo, nil
	} else {
		return UserInfo{}, ErrParseToken
	}
}

func AccountLogin(w http.ResponseWriter, r *http.Request) {

	token, err := GenerateToken(UserInfo{Id: 1})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("authorization")
	if len(auth) < 10 {
		w.Write([]byte("请先登录"))
		return
	}
	user, err := VerifyToken(auth[7:])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user.Id)
}
