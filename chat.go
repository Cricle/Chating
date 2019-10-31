package main

import (
	"Chating/chat"
	"Chating/data"
	"container/list"
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/go-xorm/xorm"
)

const (
	MSG_TYPE_SYSTEM = 0
	MSG_TYPE_USER   = 1
)

var (
	err_TokenUnauthorized = errors.New("token unauthorized!")
	succeed_Status        = &chat.StatusResponse{Status: true}
	fail_Status           = &chat.StatusResponse{Status: false}
)

type chatItem struct {
	Ser   chat.Chat_RecvServer
	Token string
	Name  string
}
type ChatService struct {
	mutex   sync.Mutex
	user    *list.List
	eng     *xorm.Engine
	userMgr *data.UserManager
}

func NewChatService(dbname string, expTimeMs time.Duration, walkMaxCount int64, walkTime time.Duration) (cs *ChatService) {
	var err error
	cs = new(ChatService)
	cs.user = list.New()
	cs.eng, err = data.NewDbContext(dbname)
	if err != nil {
		panic("db init err:" + err.Error())
	}
	cs.userMgr = data.NewUserManager(cs.eng, expTimeMs)
	cs.userMgr.LoginManager().BeginWalkDel(walkMaxCount, walkTime)
	return
}
func (cs *ChatService) Recv(req *chat.RecvRequest, ser chat.Chat_RecvServer) error {
	info, ok := cs.userMgr.LoginManager().Get(req.Token)
	if !ok {
		return err_TokenUnauthorized
	}
	cs.walk(func(c *chatItem) error {
		return c.Ser.Send(&chat.RecvResponse{
			From: info.Name,
			Pkg: &chat.SendPkg{
				Data: []byte(info.Name + " join"),
				Type: MSG_TYPE_SYSTEM,
			},
		})
	})
	cs.mutex.Lock()
	cs.user.PushBack(&chatItem{Name: info.Name, Token: req.Token, Ser: ser})
	cs.mutex.Unlock()
	<-ser.Context().Done()
	return nil
}
func (cs *ChatService) walk(f func(*chatItem) error) {
	h := cs.user.Front()
	for h != nil {
		err := f(h.Value.(*chatItem))
		if err != nil {
			fmt.Println(err.Error())
		}
		h = h.Next()
	}
}
func (cs *ChatService) Send(ctx context.Context, req *chat.SendRequest) (*chat.StatusResponse, error) {
	info, ok := cs.userMgr.LoginManager().Get(req.Token)
	if !ok {
		return nil, err_TokenUnauthorized
	}
	sendData := &chat.RecvResponse{
		From: info.Name,
		Pkg:  req.Pkg,
	}
	if req.To == "" {
		cs.walk(func(c *chatItem) error {
			return c.Ser.Send(sendData)
		})
	} else {
		cp, err := regexp.Compile(req.To)
		if err != nil {
			return fail_Status, err
		}
		cs.walk(func(c *chatItem) error {
			if cp.MatchString(c.Name) {
				return c.Ser.Send(sendData)
			}
			return nil
		})
	}
	return succeed_Status, nil
}
func (cs *ChatService) Login(ctx context.Context, req *chat.UserRequest) (cr *chat.LoginResponse, err error) {
	lr := cs.userMgr.Login(req.Name, req.Pwd)
	cr = &chat.LoginResponse{
		Status: lr.Succeed,
		Token:  lr.Token,
	}
	err = lr.Err
	return
}
func (cs *ChatService) Register(ctx context.Context, req *chat.UserRequest) (cr *chat.StatusResponse, err error) {
	rr, err := cs.userMgr.Register(req.Name, req.Pwd)
	cr = &chat.StatusResponse{
		Status: rr,
	}
	return
}
func (cs *ChatService) Logout(ctx context.Context, req *chat.LogoutRequest) (cr *chat.StatusResponse, err error) {
	cs.userMgr.Logout(req.Token)
	return succeed_Status, nil
}
