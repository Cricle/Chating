package data

import (
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func NewDbContext(dbname string) (eng *xorm.Engine, err error) {
	eng, err = xorm.NewEngine("mysql", "root:355343@/chating?charset=utf8")
	if err == nil {
		err = eng.Sync2(new(User))
	}
	return
}

type LoginRes struct {
	Token   string
	Succeed bool
	Err     error
}
type UserManager struct {
	eng      *xorm.Engine
	loginmgr *LoginManager
}

func NewUserManager(eng *xorm.Engine, expTimeMs time.Duration) *UserManager {
	um := new(UserManager)
	um.eng = eng
	um.loginmgr = NewLoginManager(expTimeMs)
	return um
}
func (um *UserManager) LoginManager() *LoginManager {
	return um.loginmgr
}
func (um *UserManager) Register(name, pwd string) (bool, error) {
	res, err := um.eng.Exec("insert into user(name,pwd,create_time) values(?,?,?)", name, pwd, time.Now())
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (um *UserManager) Login(name, pwd string) (r LoginRes) {
	res, err := um.eng.QueryString("select id from user where name=? and pwd=? limit 1", name, pwd)
	if err != nil {
		r.Err = err
		return
	}
	if len(res) == 0 {
		return
	}
	id, err := strconv.ParseInt(res[0]["id"], 10, 64)
	if err != nil {
		r.Err = err
		return
	}
	r.Token, err = um.loginmgr.New(id, name)
	if err != nil {
		r.Err = err
		return
	}
	r.Succeed = true
	return
}
func (um *UserManager) Logout(token string) {
	um.loginmgr.Del(token)
}
