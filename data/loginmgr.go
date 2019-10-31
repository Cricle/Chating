package data

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	MIN_DURATION_TIME = 10 * time.Second
)

type LoginInfo struct {
	Uid        int64
	Name       string
	CreateTime time.Time
}
type LoginManager struct {
	conc      sync.Map
	expTimeMs time.Duration
	count     int64
	walkBegin bool
}

func NewLoginManager(expTimeMs time.Duration) (lm *LoginManager) {
	if expTimeMs <= 0 {
		panic("expTimeMsg must >0")
	}
	lm = new(LoginManager)
	lm.conc = sync.Map{}
	lm.expTimeMs = expTimeMs
	return
}
func (lm *LoginManager) Count() int64 {
	return lm.count
}
func (lm *LoginManager) New(uid int64, name string) (token string, err error) {
	ut, err := uuid.NewV4()
	if err == nil {
		token = ut.String()
		info := LoginInfo{
			Uid:        uid,
			Name:       name,
			CreateTime: time.Now(),
		}
		lm.conc.Store(token, info)
		atomic.AddInt64(&lm.count, 1)
	}
	return
}
func (lm *LoginManager) Del(token string) {
	_, ok := lm.conc.Load(token)
	if ok {
		lm.conc.Delete(token)
		atomic.AddInt64(&lm.count, -1)
	}
}
func (lm *LoginManager) Get(token string) (info LoginInfo, ok bool) {
	obj, ok := lm.conc.Load(token)
	if ok {
		info = obj.(LoginInfo)
		if lm.IsOutTime(info) {
			fmt.Println("out time")
			ok = false
			lm.Del(token)
		}
	}
	return
}
func (lm *LoginManager) Clear() {
	lm.conc = sync.Map{}
	atomic.StoreInt64(&lm.count, 0)
}
func (lm *LoginManager) IsOutTime(lf LoginInfo) bool {
	return lf.CreateTime.Add(lm.expTimeMs).Unix() < time.Now().Unix()
}
func (lm *LoginManager) BeginWalkDel(maxCount int64, durationTime time.Duration) {
	if durationTime < MIN_DURATION_TIME {
		str := fmt.Sprintf("Min duration time is %vs, but now is %vs", MIN_DURATION_TIME, (durationTime / time.Second))
		panic(str)
	}
	if !lm.walkBegin {
		lm.walkBegin = true
		go func() {
			var walkCount int64
			var skipCount int64
			delKeys := make([]interface{}, 0)
			for {
				if lm.count != 0 {
					walkCount = rand.Int63() % lm.Count()
					skipCount = rand.Int63() % lm.Count()
					if walkCount > maxCount {
						walkCount = maxCount
					}
					if walkCount != 0 {
						var lf LoginInfo
						lm.conc.Range(func(key, value interface{}) bool {
							if skipCount > 0 {
								skipCount--
								return true
							}
							lf = value.(LoginInfo)
							if lm.IsOutTime(lf) {
								delKeys = append(delKeys, key)
							}
							walkCount--
							return walkCount > 0
						})
						lenDel := len(delKeys)
						for i := 0; i < lenDel; i++ {
							lm.Del(delKeys[i].(string))
						}
					}
					delKeys = make([]interface{}, 0)
				}

				<-time.After(durationTime)
			}
		}()
	}

}
