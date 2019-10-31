package main

import (
	"Chating/client/funcs"
	"Chating/client/initer"
	"bufio"
	"fmt"
	"log"
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
	fmt.Print("Token:")
	line, _, err := stdin.ReadLine()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	token := string(line)
	conn.SetToken(token)
	conn.PipManager().
		UseRecv(funcs.ImgRecv, funcs.ColorFulRev).
		UseDefaultRecv()
	color.FgGreen.Printf("初始化完成!\n已经连接到%v\n开始接收数据\n", vinit.Config().Address)
	conn.Init()
	conn.StartRev(true)
	<-exit
	conn.Close()
}
