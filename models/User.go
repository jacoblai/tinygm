package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//用户信息表
var T_User = "user"

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty" jsonschema:"-"`                                                                                        //id
	Nickname       string             `json:"nickname,omitempty" bson:"nickname,omitempty" validate:"required,min=1,max=20" jsonschema:"required,minLength=1,maxLength=20" ` //昵称
	MechanismName  string             `json:"mechanism_name,omitempty" bson:"mechanism_name,omitempty" validate:"min=0,max=20" jsonschema:"minLength=0,maxLength=20"`        //机构认证名称(主体仅有)
	Avatar         string             `json:"avatar" bson:"avatar,omitempty"`                                                                                                //头像
	Brief          string             `json:"brief,omitempty" bson:"brief,omitempty" validate:"min=0,lte=280" jsonschema:"minLength=0,maxLength=280" `                       //简介
	OpenID         string             `json:"openid,omitempty" bson:"openid,omitempty" validate:"required" jsonschema:"required"`                                            //微信用户唯一标识
	Phone          string             `json:"phone,omitempty" bson:"phone,omitempty"`                                                                                        //手机号
	Gender         string             `json:"gender,omitempty" bson:"gender,omitempty" validate:"required,oneof=男 女" jsonchema:"required,enum=男|女"`                          //男 女
	Pwd            string             `json:"pwd,omitempty" bson:"pwd,omitempty"  `                                                                                          //密码 json:"-"。在大多数操作中，都不返回给前端
	Type           string             `json:"type,omitempty" bson:"type,omitempty" validate:"required,oneof=个人 主体"  jsonchema:"required,enum=个人|主体"`                         //用户类型
	Country        string             `json:"country,omitempty" bson:"country,omitempty"`                                                                                    //国家
	Province       string             `json:"province" bson:"province,omitempty"`                                                                                            //省份
	City           string             `json:"city" bson:"city,omitempty"`                                                                                                    //城市
	FansNums       uint64             `json:"fans_nums,omitempty" bson:"fans_nums" jsonschema:"-"`                                                                           //粉丝总数
	FollowerNums   uint64             `json:"follower_nums,omitempty" bson:"follower_nums" jsonschema:"-"`                                                                   //关注总数
	VideoNum       uint64             `json:"video_num,omitempty" bson:"video_num" jsonschema:"-"`                                                                           //视频总数
	VideoTime      uint64             `json:"video_time,omitempty" bson:"video_time" jsonschema:"-"`                                                                         //视频时长，默认为30秒，以秒为单位
	Status         string             `json:"status,omitempty" bson:"status,omitempty" validate:"required,oneof=正常 冻结"jsonschema:"required,enum=正常|冻结"`                      //账户有效状态
	DisableSendMsg bool               `json:"disable_sendmsg,omitempty" bson:"disable_sendmsg,omitempty"`                                                                    //是否禁言
	DisableComment bool               `json:"disable_comment,omitempty" bson:"disable_comment,omitempty"`                                                                    //是否可以删除评论
	RegisterIp     string             `json:"register_ip,omitempty" bson:"register_ip,omitempty"`                                                                            //注册信息IP
	CreateDate     time.Time          `json:"create_date,omitempty" bson:"create_date,omitempty"`                                                                            //创建时间
	UpdateDate     time.Time          `json:"update_date,omitempty" bson:"update_date,omitempty"`                                                                            //修改时间
	SubMessagePush bool               `json:"sub_message_push,omitempty" bson:"sub_message_push,omitempty"`                                                                  //用户是否允许推送消息通知
	Images         []string           `json:"images,omitempty" bson:"images,omitempty"`                                                                                      //头像id数组(用于主体切换头像oss删除功能)
	IdentityType   string             `json:"identity_type" bson:"identity_type,omitempty" jsonchema:"enum=拥有者|运营者"`                                                         //个人用户的身份类型
	IsReview       bool               `json:"is_review" bson:"is_review" `                                                                                                   //主体是否需要审核运营者视频
	TemplateId     string             `json:"template_id,omitempty" jsonchema:"-"`                                                                                           //微信消息订阅消息模板id //数据库不存储
	Amount         int64              `json:"amount,omitempty"bson:"amount" jsonchema:"-"`                                                                                   //个人金额
}

//用户--主体关系映射表
var T_UserIdentity = "user_identity"

type UserIdentity struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty" jsonschema:"-"`                           //id
	IdentityOpenid string             `json:"identity_openid,omitempty" bson:"identity_openid,omitempty"`       //主体openid
	UserOpenid     string             `json:"user_openid,omitempty" bson:"user_openid,omitempty" `              //主体绑定的拥有者openid
	OptionOpenid   string             `json:"option_openid,omitempty" bson:"option_openid,omitempty"`           //主体绑定的运营者openid (待定)
	Status         string             `json:"status,omitempty" bson:"status,omitempty" jsonschema:"enum=正常|冻结"` //账户有效状态
	CreateDate     time.Time          `json:"create_date,omitempty" bson:"create_date,omitempty"`               //创建时间
	UpdateDate     time.Time          `json:"update_date,omitempty" bson:"update_date,omitempty"`               //修改时间
	Remark         string             `json:"remark,omitempty" bson:"remark,omitempty"`                         //备注
}

//个人中心返回数据
type UserInfo struct {
	User map[string]interface{} `json:"user" bson:"user" jsonschema:"-"` //当前用户信息
	//VideoList    []Videos               `json:"video_list" bson:"video_list" jsonschema:"-"` //视频列表
	VideoCount   uint64 `json:"video_count" bson:"video_count"`     //视频总数
	FansNums     uint64 `json:"fans_nums" bson:"fans_num" `         //粉丝总数
	FollowerNums uint64 `json:"follower_nums" bson:"follower_nums"` //专注总数
	LikeNums     uint64 `json:"like_nums" bson:"like_nums"`         //喜欢视频总数
	ReviewNums   uint64 `json:"review_nums" bson:"review_nums"`     //待审核视频总数
}
type Videos struct {
	VideoInfo map[string]interface{} `json:"video_info" bson:"video_info" jsonschema:"-"`
	CoverPath string                 `json:"cover_path" bson:"cover_path"`
}

//分享用户返回数据
type ShareUserInfo struct {
	User         map[string]interface{}   `json:"user" bson:"user" jsonschema:"-"`             //当前用户信息
	VideoList    []map[string]interface{} `json:"video_list" bson:"video_list" jsonschema:"-"` //视频列表
	VideoCount   uint64                   `json:"video_count" bson:"video_count"`              //视频总数
	FansNums     uint64                   `json:"fans_nums" bson:"fans_num" `                  //粉丝总数
	FollowerNums uint64                   `json:"follower_nums" bson:"attention_num"`          //专注总数
	IsFollowed   bool                     `json:"is_followed" bson:"is_followed"`              //是否关注
}
