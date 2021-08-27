package engine

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/jacoblai/mschema"
	"github.com/jacoblai/tinygm/models"
	"github.com/jacoblai/yiyidb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"runtime"
	"time"
)

var jwtSign = []byte("Tc2020DseslT_Mgt!@#")
var (
	WxAppId               = "wxf00cf4bcda0a98e8"                          // 微信appid
	WxSecret              = "e1d2c22e3361849299bbf46c6141b90b"            // AppSecret
	WxMchID               = "1584216311"                                  // 商户号
	WxApiKey              = "66e20643t3f0ft43f3tabfat5dab6340"            // api key v3
	WxSerialNo            = "50C31465C66676EF4AB0B4B291E07CF15ECA38CC"    //商户号证书序列号
	NotifyBaseUrl         = "https://umall.yiyii.net"                     //本服务 公网域名
	WxTemplateId          = "rEMgZ3t_jL5lU0zjSYbtrPOeBn9cmVueCjPFOzlfSQg" //小程序订阅消息模板id
	WxTemplateThing1      = "thing4"                                      //小程序订阅消息值1
	WxTemplateThing2      = "thing12"
	WxMchReviewTemplateId = "Z_iVCFSeWhL3BqnL_E7o7kVRc-z6N79Isbsviko0lnk" // //商家入驻 小程序订阅消息模板id
)

type DbEngine struct {
	MgEngine *mongo.Client //数据库引擎
	validate *validator.Validate
	TTlDb    *yiyidb.TtlRunner
	TokenDb  *yiyidb.Kvdb //token引擎
	Mdb      string
}

func NewDbEngine() *DbEngine {
	return &DbEngine{
		validate: validator.New(),
	}
}

func (d *DbEngine) Open(dir, mg, mdb string) error {
	d.Mdb = mdb

	ttldb, err := yiyidb.OpenTtlRunner(dir + "/ttl.db")
	if err != nil {
		return err
	}
	d.TTlDb = ttldb

	tokenDb, err := yiyidb.OpenKvdb(dir+"/token.db", 24)
	if err != nil {
		return err
	}
	err = tokenDb.RegTTl("token.db", d.TTlDb)
	if err != nil {
		return err
	}
	d.TokenDb = tokenDb

	ops := options.Client().ApplyURI(mg)
	p := uint64(runtime.NumCPU() * 2)
	ops.MaxPoolSize = &p
	ops.WriteConcern = writeconcern.New(writeconcern.J(true), writeconcern.W(1))
	ops.ReadPreference = readpref.PrimaryPreferred()
	db, err := mongo.NewClient(ops)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = db.Connect(ctx)
	if err != nil {
		return err
	}
	err = db.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		log.Printf("ping err:%v", err)
	}

	d.MgEngine = db

	var session *mongo.Client
	ss, err := mongo.NewClient(ops)
	if err != nil {
		panic(err)
	}
	err = ss.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	session = ss
	defer session.Disconnect(context.Background())

	//用户表
	res := InitDbAndColl(session, mdb, models.T_User, GenJsonSchema(&models.User{}))
	users := session.Database(mdb).Collection(models.T_User)
	indexView := users.Indexes()
	_, err = indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{"nickname", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"phone", bsonx.Int32(1)}},
		},
		{
			Keys:    bsonx.Doc{{"openid", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Println(err)
	}
	log.Println(models.T_User, res["ok"])
	return nil
}

func (d *DbEngine) GetColl(coll string) *mongo.Collection {
	col, _ := d.MgEngine.Database(d.Mdb).Collection(coll).Clone()
	return col
}

func InitDbAndColl(session *mongo.Client, db, coll string, model map[string]interface{}) map[string]interface{} {
	tn, _ := session.Database(db).ListCollections(context.Background(), bson.M{"name": coll})
	if tn.Next(context.Background()) == false {
		session.Database(db).RunCommand(context.Background(), bson.D{{"create", coll}})
	}
	result := session.Database(db).RunCommand(context.Background(), bson.D{{"collMod", coll}, {"validator", model}})
	var res map[string]interface{}
	err := result.Decode(&res)
	if err != nil {
		log.Println(err)
	}
	return res
}

//创建数据库验证schema结构对象

func GenJsonSchema(obj interface{}) map[string]interface{} {
	flect := &mschema.Reflector{ExpandedStruct: true, RequiredFromJSONSchemaTags: true, AllowAdditionalProperties: true}
	ob := flect.Reflect(obj)
	bts, _ := json.Marshal(&ob)
	var o map[string]interface{}
	_ = json.Unmarshal(bts, &o)
	return bson.M{"$jsonSchema": o}
}

func (d *DbEngine) Close() {
	_ = d.MgEngine.Disconnect(context.Background())
}
