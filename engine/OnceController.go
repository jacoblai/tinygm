package engine

import (
	"github.com/jacoblai/httprouter"
	"github.com/jacoblai/tinygm/resultor"
	"github.com/jacoblai/yiyidb"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
	"sync"
)

//取防重复因子
func (d *DbEngine) GetOnce(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	var lock *sync.Mutex
	if lk, ok := d.ipDenyLocks.Load(ip); ok {
		lock = lk.(*sync.Mutex)
	} else {
		lock = &sync.Mutex{}
		d.ipDenyLocks.Store(ip, lock)
	}
	has := d.ipDenyDb.Exists([]byte(ip), nil)
	if has {
		v, _ := d.ipDenyDb.Get([]byte(ip), nil)
		ct := yiyidb.KeyToIDPure(v)
		if ct > 1000 {
			resultor.RetErr(w, "ip已到上限")
			return
		}
		lock.Lock()
		ct++
		lock.Unlock()
		_ = d.ipDenyDb.Put([]byte(ip), yiyidb.IdToKeyPure(ct), 0, nil)
	} else {
		_ = d.ipDenyDb.Put([]byte(ip), yiyidb.IdToKeyPure(1), 3600, nil)
	}
	onceId := uuid.NewV4().String()
	once := []byte(onceId)
	err := d.CacheDb.Put(once, once, 36000, nil)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	resultor.RetOk(w, onceId, 1)
}
