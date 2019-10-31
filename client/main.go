package main

import (
	"Chating/client/funcs"
	"Chating/client/initer"
	"Chating/client/tiping"
	"bufio"
	"fmt"
	"os"

	"github.com/gookit/color"
)

var (
	vinit = initer.NewIniter()
	conn  *initer.Connector
	exit  = make(chan bool, 0)
)

func main() {

	stdin := bufio.NewReader(os.Stdin)
	vinit.Init(stdin)

	conn, err := initer.NewConnector(vinit.Config())
	if err != nil {
		fmt.Println("did not connect: %v", err)
		return
	}
	conn.PipManager().
		UseRecv(funcs.ImgRecv, funcs.ColorFulRev).
		UseDefaultRecv()
	conn.PipManager().
		UseSend(funcs.WatchFile, funcs.ImgSend, funcs.ColorFulSend, funcs.TipSend).
		UseSendTo([]byte("to")).
		UseSend(conn.LoginHandle(stdin, vinit, false)).
		UseSend(conn.RegisterHandle(stdin, vinit)).
		UseSend(conn.TokenHandle()).
		UseSend(QuitHandle)
	tiping.DefaultTiper.Add("quit", tiping.MakeDefaultPrint(color.FgDefault, "quit", "退出"))
	conn.EndHandleSend()
	color.FgGreen.Printf("初始化完成!\n已经连接到%v\n/help查看可使用的命令\n", vinit.Config().Address)
	conn.Begin(stdin)
	<-exit
	conn.Close()
}
func QuitHandle(sc *initer.SendContext) {
	sc.DefaultHandleCommand(func(comm string, args []string) bool {
		if comm == "quit" {
			if conn != nil && conn.ChatRevClient() != nil {
				conn.Close()
			}
			os.Exit(0)
			return true
		}
		return false
	})
}
