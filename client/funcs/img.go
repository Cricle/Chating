package funcs

import (
	"Chating/client/initer"
	"Chating/client/tiping"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gookit/color"
	uuid "github.com/satori/go.uuid"
)

const (
	COMMAND_IMG   = "img"
	COMMAND_WATCH = "watch"

	MEDATA_FILENAME = "filename"

	TYPE_IMG = 3
)

var (
	mutex                       = sync.RWMutex{}
	imgMapper map[uint64]string = make(map[uint64]string, 0)
	idx                         = uint64(0)
)

func init() {
	tiping.DefaultTiper.Add(COMMAND_IMG, tiping.MakeDefaultPrint(color.FgDefault, COMMAND_IMG, "发送图片", tiping.CmdArg{Name: "path", Descript: "本地图片路径"}))
}
func ImgRecv(rc *initer.RecvContext) {
	if rc.EqualType(TYPE_IMG) {
		fn := rc.Response().Pkg.Medata[MEDATA_FILENAME]
		if fn == nil {
			fmt.Println("Error file name empty")
		} else {
			ud, _ := uuid.NewV4()
			strfn := string(fn)
			savefilename := ud.String() + filepath.Ext(strfn)
			mutex.Lock()
			curridx := idx
			imgMapper[idx] = savefilename
			idx++
			mutex.Unlock()
			err := ioutil.WriteFile(savefilename, rc.Response().Pkg.Data, os.ModePerm)
			if err != nil {
				fmt.Println("Create file " + strfn + " error")
				return
			}
			if rc.Appear {
				np := color.FgGray
				ip := color.FgLightCyan
				np.Print("Recv img name is ")
				ip.Printf("%s", strfn)
				np.Print("(Save as ")
				ip.Printf("%s", savefilename)
				np.Print(") -> id=")
				ip.Printf("%v\n", curridx)
			}

		}
		rc.Done()
	}
}
func ImgSend(rc *initer.SendContext) {
	comm, args := rc.DefaultCommandStringSplit()
	if comm == COMMAND_IMG {
		bs, err := ioutil.ReadFile(args)
		if err != nil {
			rc.Error(err)
		} else {
			rc.Data(bs)
			rc.PkgType(TYPE_IMG)
			fns := strings.Split(args, "\\")
			rc.Model().Pkg.Medata[MEDATA_FILENAME] = []byte(fns[len(fns)-1])
			rc.Done()
		}
	}
}
func show(args ...string) {
	cd := exec.Command("explorer.exe", args...)
	cd.Run()
}
func WatchFile(rc *initer.SendContext) {
	comm, args := rc.DefaultCommandStringSplit()
	if comm == COMMAND_WATCH {
		index, err := strconv.ParseUint(args, 10, 64)
		if err == nil {
			mutex.RLock()
			fn := imgMapper[index]
			mutex.RUnlock()
			if fn != "" {
				show(fn)
			} else {
				show(args)
			}
		} else {
			show(args)
		}
		rc.Done()
		rc.SetHandle(true)
	}
}
