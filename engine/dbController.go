package engine

import (
	"context"
	"github.com/jacoblai/httprouter"
	"github.com/jacoblai/tinygm/deny"
	"github.com/jacoblai/tinygm/resultor"
	"github.com/pquerna/ffjson/ffjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
)

//新增

func (d *DbEngine) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	//防重复提交
	once := r.URL.Query().Get("once")
	if once == "" {
		resultor.RetErr(w, "空参数防重放因子")
		return
	}
	_, err := d.CacheDb.Get([]byte(once), nil)
	if err != nil {
		resultor.RetErr(w, "防重放或因子过期生效，请重试")
		return
	} else {
		_ = d.CacheDb.Del([]byte(once), nil)
	}

	collName := r.URL.Query().Get("cn")
	if collName == "" {
		resultor.RetErr(w, "集合名为空")
		return
	}

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

	var obj map[string]interface{}
	err = ffjson.Unmarshal(body, &obj)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}

	id, err := d.GetColl(collName).InsertOne(context.Background(), &obj)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}

	resultor.RetOk(w, map[string]interface{}{"id": id.InsertedID}, 1)
}

//修改

func (d *DbEngine) Put(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	collName := r.URL.Query().Get("cn")
	if collName == "" {
		resultor.RetErr(w, "集合名为空")
		return
	}

	var obj map[string]interface{}
	err = ffjson.Unmarshal(body, &obj)
	if err != nil {
		resultor.RetErr(w, "解析参数出错")
		return
	}

	filter, ok := obj["filter"].(map[string]interface{})
	if !ok {
		resultor.RetErr(w, "filter invalidate")
		return
	}

	update, ok := obj["update"].(map[string]interface{})
	if !ok {
		resultor.RetErr(w, "update invalidate")
		return
	}

	res, err := d.GetColl(collName).UpdateMany(context.Background(), filter, update)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	resultor.RetOk(w, res, int(res.ModifiedCount))
}

//部分修改

func (d *DbEngine) Patch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	collName := r.URL.Query().Get("cn")
	if collName == "" {
		resultor.RetErr(w, "集合名为空")
		return
	}

	var obj map[string]interface{}
	err = ffjson.Unmarshal(body, &obj)
	if err != nil {
		resultor.RetErr(w, "解析参数出错")
		return
	}

	filter, ok := obj["filter"].(map[string]interface{})
	if !ok {
		resultor.RetErr(w, "filter invalidate")
		return
	}

	update, ok := obj["update"].(map[string]interface{})
	if !ok {
		resultor.RetErr(w, "update invalidate")
		return
	}

	res, err := d.GetColl(collName).UpdateMany(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	resultor.RetOk(w, res, int(res.ModifiedCount))
}

//获取

func (d *DbEngine) Find(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	collName := r.URL.Query().Get("cn")
	if collName == "" {
		resultor.RetErr(w, "集合名为空")
		return
	}

	var query []map[string]interface{}
	err = ffjson.Unmarshal(body, &query)
	if err != nil {
		resultor.RetErr(w, "解析参数出错")
		return
	}

	data, err := d.GetColl(collName).Aggregate(context.Background(), query)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	res := make([]map[string]interface{}, 0)
	err = data.All(context.Background(), &res)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	resultor.RetOk(w, res, len(res))
}

//获取某个

func (d *DbEngine) FindById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		resultor.RetErr(w, "1003")
		return
	}

	collName := r.URL.Query().Get("cn")
	if collName == "" {
		resultor.RetErr(w, "集合名为空")
		return
	}
	data := d.GetColl(collName).FindOne(context.Background(), bson.M{"_id": id})
	if data.Err() != nil {
		resultor.RetErr(w, data.Err().Error())
		return
	}
	var result map[string]interface{}
	err = data.Decode(&result)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	resultor.RetOk(w, result, 1)
}

//删除

func (d *DbEngine) Del(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	collName := r.URL.Query().Get("cn")
	if collName == "" {
		resultor.RetErr(w, "集合名为空")
		return
	}

	var obj map[string]interface{}
	err = ffjson.Unmarshal(body, &obj)
	if err != nil {
		resultor.RetErr(w, "解析参数出错")
		return
	}

	filter, ok := obj["filter"].(map[string]interface{})
	if !ok {
		resultor.RetErr(w, "filter invalidate")
		return
	}

	res, err := d.GetColl(collName).DeleteMany(context.Background(), filter)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	resultor.RetOk(w, res, int(res.DeletedCount))
}
