package initer

import (
	"Chating/chat"
	"Chating/client/tiping"
	"Chating/helper"
	"bufio"
	"context"
	"errors"
	"fmt"
	"time"

	color "github.com/gookit/color"
	"google.golang.org/grpc"
)

const (
	COMMAND_LOGIN    = "login"
	COMMAND_REGISTER = "register"
	COMMAND_TOKEN    = "token"
)

func init() {
	tiping.DefaultTiper.Add(COMMAND_LOGIN, tiping.MakeDefaultPrint(color.FgDefault, COMMAND_LOGIN, "登录"))
	tiping.DefaultTiper.Add(COMMAND_REGISTER, tiping.MakeDefaultPrint(color.FgDefault, COMMAND_REGISTER, "注册"))
}

var (
	Error_NoAuthorized = errors.New("Has not logined")
	Error_LoginFail    = errors.New("Login fail")
	Error_RegisterFail = errors.New("Register fail")
)

type Connector struct {
	conn         *grpc.ClientConn
	pm           *PipManager
	client       chat.ChatClient
	revClient    chat.Chat_RecvClient
	config       *Config
	token        string
	tokenExpTime time.Time
	inited       bool
	appearMsg    bool
	endSend      bool
}

func NewConnector(config *Config) (c *Connector, err error) {
	c = new(Connector)
	c.config = config
	c.pm = NewPipManager()
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	c.conn = conn
	c.client = chat.NewChatClient(c.conn)
	c.appearMsg = true
	return c, err
}
func (c *Connector) SetToken(token string) {
	c.token = token
}
func (c *Connector) PipManager() *PipManager {
	return c.pm
}
func (c *Connector) RevClient() chat.ChatClient {
	return c.client
}
func (c *Connector) ChatRevClient() chat.Chat_RecvClient {
	return c.revClient
}
func (c *Connector) IsInited() bool {
	return c.inited
}
func (c *Connector) Login(name, pwd string) error {
	md5Pwd, err := helper.MakeMd5Str(pwd)
	if err != nil {
		return err
	}
	rep, err := c.client.Login(context.Background(), &chat.UserRequest{
		Name: name,
		Pwd:  md5Pwd,
	})
	if err != nil {
		return err
	}
	if !rep.Status {
		return Error_LoginFail
	}
	c.tokenExpTime = time.Now().Add(time.Duration(rep.ExpTime))
	c.token = rep.Token
	return err
}
func (c *Connector) Register(name, pwd string) error {
	md5Pwd, err := helper.MakeMd5Str(pwd)
	if err != nil {
		return err
	}
	rep, err := c.client.Register(context.Background(), &chat.UserRequest{
		Name: name,
		Pwd:  md5Pwd,
	})
	if err != nil {
		return err
	}
	if !rep.Status {
		err = Error_RegisterFail
	}
	return err
}
func (c *Connector) Init() error {
	if c.token == "" {
		return Error_NoAuthorized
	}
	rev, err := c.client.Recv(context.Background(), &chat.RecvRequest{Token: c.token})
	if err != nil {
		return err
	}
	c.revClient = rev
	if err == nil {
		c.inited = true
	}
	return err
}
func (c *Connector) StartRev(appear bool) error {
	if c.token == "" {
		return Error_NoAuthorized
	}
	c.appearMsg = appear
	go func() {
		for {
			rep, err := c.revClient.Recv()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			c.pm.HandleRecv(rep, c.appearMsg)
		}
	}()
	return nil
}
func (c *Connector) Begin(stdin *bufio.Reader) error {
	for {
		if c.PipManager().GetReadyHandle() != nil {
			c.PipManager().GetReadyHandle()()
		}
		line, _, err := stdin.ReadLine()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if len(line) == 0 {
			continue
		}
		sc := c.pm.HandleSend(line)
		if !sc.handle && c.endSend {
			if c.token == "" {
				sc.Error(Error_NoAuthorized)
			} else {
				if !sc.done {
					sc.Data(sc.Line())
					sc.PkgType(1)
				}
				sc.model.Token = c.token
				_, err := c.client.Send(context.Background(), sc.model)
				if err != nil {
					sc.Error(err)
				}
				sc.Done()
				sc.SetHandle(true)
			}
		}
		if !sc.GetHandle() {
			color.FgYellow.Println("WARN:Command has not handled")
		}
		if sc.err != nil {
			color.FgRed.Println(sc.err.Error())
		}
	}
	return nil
}
func (c *Connector) EndHandleSend() {
	c.endSend = true

}
func (c *Connector) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
func (c *Connector) LoginHandle(in *bufio.Reader, i *Initer, appearMsg bool) func(sp *SendContext) {
	return func(sp *SendContext) {
		sp.DefaultHandleCommand(func(comm string, args []string) bool {
			if comm == COMMAND_LOGIN {
				i.WithConfig(in)
				cfg := i.Config()
				err := c.Login(cfg.Name, cfg.Pwd)
				if err != nil {
					sp.Error(err)
					return true
				} else {
					color.FgGreen.Println("Login succeed")
				}
				err = i.Init(in)
				if err != nil {
					sp.Error(err)
					return true
				}
				err = c.Init()
				if err != nil {
					sp.Error(err)
					return true
				}
				err = c.StartRev(appearMsg)
				if err != nil {
					sp.Error(err)
					return true
				}
				return true
			}
			return false
		})
	}
}
func (c *Connector) TokenHandle() func(sp *SendContext) {
	return func(sp *SendContext) {
		sp.DefaultHandleCommand(func(comm string, args []string) bool {
			if comm == COMMAND_TOKEN {
				fmt.Println(c.token)
				return true
			}
			return false
		})
	}
}
func (c *Connector) RegisterHandle(in *bufio.Reader, i *Initer) func(sp *SendContext) {
	return func(sp *SendContext) {
		sp.DefaultHandleCommand(func(comm string, args []string) bool {
			if comm == COMMAND_REGISTER {
				i.WithConfig(in)
				cfg := i.Config()
				err := c.Register(cfg.Name, cfg.Pwd)
				if err != nil {
					sp.Error(err)
				}
				return true
			}
			return false
		})
	}
}
