package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserId         string `json:"user_id"`         //当前登陆者id
	UserType       string `json:"user_type"`       //主体、个人
	UserOpenId     string `json:"user_openid"`     //当前登陆者的openid
	UserStatus     string `json:"user_status"`     //当前登录者状态
	OperatorOpenID string `json:"operator_openid"` //主体绑定的个人的openid
	OperatorType   string `json:"operator_type"`   //当前主体的类型 拥有者/运营者
	SysRole        string `json:"sys_role"`        //后台管理员身份
}

type JwtClaims struct {
	UserId  string `json:"user_id"`
	SysRole string `json:"sys_role"`
	jwt.StandardClaims
}
