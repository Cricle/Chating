package funcs

import (
	"Chating/client/initer"
	"Chating/client/tiping"
	"bytes"
	"errors"

	color "github.com/gookit/color"
)

const (
	TYPE_COLORFUL    = 4
	COMMAND_COLORFUL = "cf"
	MEDATA_COLORKEY  = "colorkey"
)

func init() {
	tiping.DefaultTiper.Add(COMMAND_COLORFUL, tiping.MakeDefaultPrint(color.FgDefault, COMMAND_COLORFUL, "发送彩色字",
		tiping.CmdArg{Name: "colorKey", Descript: "颜色键值:[red,blue,green,yellow,default,black,magenta,cyan,white]"},
		tiping.CmdArg{Name: "...text", Descript: "发送的文本"}))
}
func ColorFulRev(rc *initer.RecvContext) {
	if rc.Response().Pkg.Type == TYPE_COLORFUL {
		col := rc.Response().Pkg.Medata[MEDATA_COLORKEY]
		if col != nil && rc.Appear {
			c := color.FgColors[string(col)]
			c.Println(string(rc.Response().Pkg.Data))
			rc.Done()
		}
	}
}
func ColorFulSend(rc *initer.SendContext) {
	comm, args := rc.DefaultCommandSplit()
	if bytes.Equal(comm, []byte(COMMAND_COLORFUL)) {
		if len(args) > 1 {
			rc.Model().Pkg.Medata[MEDATA_COLORKEY] = args[0]
			rc.Data(bytes.Join(args[1:], initer.SpliteTag))
			rc.PkgType(TYPE_COLORFUL)
			rc.Done()
		} else {
			rc.Error(errors.New("par less"))
			rc.SetHandle(true)
		}
	}
}
