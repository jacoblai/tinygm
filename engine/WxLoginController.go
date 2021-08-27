package engine

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jacoblai/httprouter"
	"github.com/jacoblai/tinygm/deny"
	"github.com/jacoblai/tinygm/models"
	"github.com/jacoblai/tinygm/resultor"
	"github.com/medivhzhan/weapp/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//判断库里是否有这个用户，没有的话就返回给前端需要授权

func (d *DbEngine) LoginOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}
	if len(body) == 0 {
		resultor.RetErr(w, "1001")
		return
	}
	//防注入
	if !deny.InjectionPass(body) {
		resultor.RetErr(w, "1002")
		return
	}
	var alb map[string]interface{}
	err = json.Unmarshal(body, &alb)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}
	//请求微信接口获取数据
	userInfo, err := d.Login(alb["code"].(string))
	if userInfo == nil || err != nil {
		//目前还没有针对请求为空应当如何处理数据的需求
		resultor.RetErr(w, "微信官方登录未获取到信息")
		return
	}
	//用拿到的数据userInfo.OpenID微信ID查数据库
	var u map[string]interface{}
	tUser := d.GetColl(models.T_User)
	uCount, err := tUser.CountDocuments(context.Background(), bson.M{"openid": userInfo.OpenID})
	//如果即没有查到数据，则返回给前端需要进行授权
	if err != nil || uCount <= 0 {
		resultor.RetOk(w, "请先进行授权", 1)
		return
	}
	//如果有这个用户就根据openID查出user表里的用户
	getUser := tUser.FindOne(context.Background(), bson.M{"openid": userInfo.OpenID})
	if getUser.Err() != nil {
		resultor.RetErr(w, getUser.Err())
		return
	}
	err = getUser.Decode(&u)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}

	oid, ok := u["_id"].(primitive.ObjectID)
	if !ok {
		resultor.RetErr(w, "内部错误")
		return
	}

	reUser := make(map[string]interface{})
	reUser["id"] = oid.Hex()
	reUser["nickname"] = u["nickname"]
	reUser["gender"] = u["gender"]
	reUser["status"] = u["status"]
	reUser["avatar"] = u["avatar"]
	reUser["type"] = u["type"]
	reUser["city"] = u["city"]
	reUser["province"] = u["province"]
	reUser["phone"] = u["phone"]
	reUser["video_time"] = u["video_time"]
	reUser["disable_comment"] = u["disable_comment"]
	//reUser["openid"] = userInfo.OpenID
	reUser["template_id"] = WxTemplateId
	if u["type"] == "主体" {
		reUser["mechanism_name"] = u["mechanism_name"]
	}
	if ot, ok := u["identity_type"]; ok {
		//创建token
		sessionKey := primitive.NewObjectID().Hex()
		err = d.TokenDb.PutObject([]byte(sessionKey), models.Claims{
			UserId:         oid.Hex(),
			UserType:       u["type"].(string),
			UserOpenId:     u["openid"].(string),
			UserStatus:     u["status"].(string),
			OperatorOpenID: userInfo.OpenID,
			OperatorType:   ot.(string),
		}, 6*3600, nil)
		if err != nil {
			resultor.RetErr(w, "token new error")
			return
		}
		//最后返回的信息
		result := make(map[string]interface{})
		result["token"] = sessionKey
		result["userInfo"] = reUser
		resultor.RetOk(w, result, 1)
	} else {
		//创建token
		sessionKey := primitive.NewObjectID().Hex()
		err = d.TokenDb.PutObject([]byte(sessionKey), models.Claims{
			UserId:         oid.Hex(),
			UserType:       u["type"].(string),
			UserOpenId:     u["openid"].(string),
			UserStatus:     u["status"].(string),
			OperatorOpenID: userInfo.OpenID,
		}, 6*3600, nil)
		if err != nil {
			resultor.RetErr(w, "token new error")
			return
		}
		//最后返回的信息
		result := make(map[string]interface{})
		result["token"] = sessionKey
		result["userInfo"] = reUser
		resultor.RetOk(w, result, 1)
	}
}

//判断库里是否有这个用户，没有的话就返回给前端需要授权

func (d *DbEngine) LoginCodeToOpenId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}
	if len(body) == 0 {
		resultor.RetErr(w, "1001")
		return
	}
	//防注入
	if !deny.InjectionPass(body) {
		resultor.RetErr(w, "1002")
		return
	}
	var alb map[string]interface{}
	err = json.Unmarshal(body, &alb)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}
	//请求微信接口获取数据
	userInfo, err := d.Login(alb["code"].(string))
	if userInfo == nil || err != nil {
		//目前还没有针对请求为空应当如何处理数据的需求
		resultor.RetErr(w, "微信官方登录未获取到信息")
		return
	}
	//最后返回的信息
	result := make(map[string]interface{})
	result["openid"] = userInfo.OpenID
	resultor.RetOk(w, result, 1)
}

func (d *DbEngine) LoginByWxAndGetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}
	if len(body) == 0 {
		resultor.RetErr(w, "1001")
		return
	}
	//防注入
	if !deny.InjectionPass(body) {
		resultor.RetErr(w, "1002")
		return
	}
	var alb models.AuthLoginBody
	err = json.Unmarshal(body, &alb)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}
	//获取ip~以冒号拆分分割IP和端口号
	clientIP := strings.Split(r.RemoteAddr, ":")
	//请求微信接口获取数据
	userInfo, err := d.LoginGetUserInfo(alb.Code, alb.UserInfo)
	if err != nil {
		//目前还没有针对请求为空应当如何处理数据的需求
		resultor.RetErr(w, "没有获取到微信用户数据"+err.Error())
		return
	}
	//微信返回gender 1,2
	var gender string
	if userInfo.Gender == 1 {
		gender = "男"
	} else {
		gender = "女"
	}
	//用拿到的数据userInfo.Openid查数据库
	var u models.User
	tUser := d.GetColl(models.T_User)
	uCount, _ := tUser.CountDocuments(context.Background(), bson.M{"openid": userInfo.OpenID})
	//如果即没有查到数据，则进行一次插入数据操作
	if uCount <= 0 {
		//默认普通用户,随机密码 用户角色层级？
		temp := md5.Sum([]byte(strconv.Itoa(rand.Intn(99999999-10000000) + 10000000)))
		newUser := models.User{
			Nickname:       userInfo.Nickname,
			Gender:         gender,
			Pwd:            hex.EncodeToString(temp[:]),
			Status:         "正常",
			Type:           "个人",
			OpenID:         userInfo.OpenID,
			Brief:          "",
			Phone:          "",
			Country:        userInfo.Country,
			Province:       userInfo.Province,
			City:           userInfo.City,
			FansNums:       0,
			FollowerNums:   0,
			VideoNum:       0,
			VideoTime:      36,
			Amount:         0,
			DisableSendMsg: false,
			DisableComment: false,
			CreateDate:     time.Now().Local(),
			UpdateDate:     time.Now().Local(),
			Avatar:         userInfo.Avatar,
			RegisterIp:     clientIP[0],
			SubMessagePush: false,
		}
		err = d.validate.Struct(newUser)
		if err != nil {
			resultor.RetErr(w, "validate err："+err.Error())
			return
		}
		//插入一条数据
		_, err := tUser.InsertOne(context.Background(), &newUser)
		if err != nil {
			resultor.RetErr(w, err.Error())
			return
		}
	}
	//不管是新用户还是老用户 根据openID查出user表里的用户
	getUser := tUser.FindOne(context.Background(), bson.M{"openid": userInfo.OpenID})
	if getUser.Err() != nil {
		resultor.RetErr(w, getUser.Err())
		return
	}
	err = getUser.Decode(&u)
	if err != nil {
		resultor.RetErr(w, err)
		return
	}
	reUser := make(map[string]interface{})
	reUser["id"] = u.ID
	reUser["nickname"] = u.Nickname
	reUser["gender"] = gender
	reUser["status"] = u.Status
	reUser["avatar"] = u.Avatar
	reUser["city"] = u.City
	reUser["province"] = u.Province
	reUser["type"] = u.Type
	reUser["video_time"] = u.VideoTime
	reUser["disable_comment"] = u.DisableComment
	//reUser["mechanism_name"] = u.MechanismName
	reUser["template_id"] = WxTemplateId
	//创建token
	sessionKey := primitive.NewObjectID().Hex()
	err = d.TokenDb.PutObject([]byte(sessionKey), models.Claims{
		UserId:         u.ID.Hex(),
		UserType:       u.Type,
		UserOpenId:     u.OpenID,
		UserStatus:     u.Status,
		OperatorOpenID: userInfo.OpenID,
	}, 6*3600, nil)
	if err != nil {
		resultor.RetErr(w, "token new error")
		return
	}
	//最后返回的信息
	result := make(map[string]interface{})
	result["token"] = sessionKey
	result["userInfo"] = reUser
	resultor.RetOk(w, result, 1)
}

func (d *DbEngine) LoginGetUserInfo(code string, uInfo models.ResUserInfo) (*weapp.UserInfo, error) {
	//登录凭证校验 获取微信用户openid，session_key,unionid
	lres, err := d.Login(code)
	if err != nil {
		return nil, err
	}
	//解密用户信息数据
	res, err := weapp.DecryptUserInfo(lres.SessionKey, uInfo.RawData, uInfo.EncryptedData, uInfo.Signature, uInfo.IV)
	if err != nil {
		return nil, err
	}
	res.OpenID = lres.OpenID
	return res, nil
}

func (d *DbEngine) Login(code string) (*weapp.LoginResponse, error) {
	//登录凭证校验 获取微信用户openid，session_key,unionid
	res, err := weapp.Login(WxAppId, WxSecret, code)
	if err != nil {
		return nil, err
	}
	if err := res.GetResponseError(); err != nil {
		return nil, err
	}
	return res, nil
}

func (d *DbEngine) WxGetUserPhone(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	if len(body) == 0 {
		resultor.RetErr(w, "1001")
		return
	}
	//防注入
	if !deny.InjectionPass(body) {
		resultor.RetErr(w, "1002")
		return
	}
	var obj map[string]string
	err = json.Unmarshal(body, &obj)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	code, ok := obj["code"]
	if !ok {
		resultor.RetErr(w, errors.New("code nil").Error())
		return
	}
	ed, ok := obj["encryptedData"]
	if !ok {
		resultor.RetErr(w, errors.New("encryptedData nil").Error())
		return
	}
	iv, ok := obj["iv"]
	if !ok {
		resultor.RetErr(w, errors.New("iv nil").Error())
		return
	}
	uinfo, err := d.Login(code)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	res, err := weapp.DecryptMobile(uinfo.SessionKey, ed, iv)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	resultor.RetOk(w, res, 1)
}
