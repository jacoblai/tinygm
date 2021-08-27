package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/jacoblai/httprouter"
	"github.com/jacoblai/tinygm/cors"
	"github.com/jacoblai/tinygm/engine"
	"github.com/jacoblai/tinygm/limit"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	var (
		addr  = flag.String("l", ":8000", "绑定Host地址")
		mongo = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		//mongo = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		db = flag.String("db", "tinygm", "database name") //数据库
	)
	flag.Parse()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(dir + "/data"); err != nil {
		log.Println(err)
	}

	//启动文件日志
	//logFile, logErr := os.OpenFile(dir+"/dal.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	//if logErr != nil {
	//	log.Printf("err: %v\n", logErr)
	//	return
	//}
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	eng := engine.NewDbEngine()
	if *addr == ":443" {
		*mongo = "mongodb://root:33a198d7-a2d1-4be7-8156-763e028816c8@172.18.157.150:56431,172.18.157.149:56431,172.18.157.148:56431/?authSource=admin&replicaSet=rs1"
		err = eng.Open(dir, *mongo, *db)
		if err != nil {
			log.Fatal("database connect error")
		}
	} else {
		err = eng.Open(dir, *mongo, *db)
		if err != nil {
			log.Fatal("database connect error")
		}
	}

	router := httprouter.New()
	rootPath := "/api"

	//登陆相关
	router.POST(rootPath+"/wx/loginOne", eng.LoginOne)                        //登录验证用户是否授权
	router.POST(rootPath+"/wx/codeToOpenId", eng.LoginCodeToOpenId)           //无感登陆
	router.POST(rootPath+"/wx/loginByWxInfo", eng.LoginByWxAndGetUserInfo)    // 登录
	router.POST(rootPath+"/wx/getWxPhone", eng.TokenAuth(eng.WxGetUserPhone)) // 解微信手机号

	srv := &http.Server{Handler: limit.Limit(cors.CORS(router)), ErrorLog: nil}
	srv.Addr = *addr
	if *addr == ":443" {
		//生产环境取消刷库后以TLS模式运行服务
		cert, err := tls.LoadX509KeyPair(dir+"/data/umall.yiyii.net.pem", dir+"/data/umall.yiyii.net.key")
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		srv.TLSConfig = config
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				log.Println(err)
			}
		}()
		log.Println("server on https port", srv.Addr)
	} else {
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}()
		log.Println("server on http port", srv.Addr)
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
			go func() {
				_ = srv.Shutdown(ctx)
				cleanup <- true
			}()
			<-cleanup
			eng.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
