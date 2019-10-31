package initer

import (
	"Chating/chat"
	"bytes"
	"errors"
	"fmt"
)

var (
	CommandPrefx = []byte("/")
	SpliteTag    = []byte(" ")
)

type RecvPipHandle func(rc *RecvContext)
type SendPipHandle func(sc *SendContext)
type ReadyHandle func()
type SendLine []byte

type SendContext struct {
	line     SendLine
	model    *chat.SendRequest
	comm     string
	arg      string
	commByte []byte
	argByte  [][]byte

	doneDefault     bool
	doneDefaultByte bool

	err    error
	done   bool
	handle bool
}

func (sc *SendContext) GetHandle() bool {
	return sc.handle
}
func (sc *SendContext) SetHandle(handle bool) {
	sc.handle = handle
}
func (sc *SendContext) PrefxSplit(prefx []byte, commTag []byte) [][]byte {
	if bytes.HasPrefix(sc.line, prefx) {
		return bytes.Split(sc.line, commTag)
	}
	return nil
}
func (sc *SendContext) CommandSplit(prefx []byte, commTag []byte) (comm []byte, par [][]byte) {
	res := sc.PrefxSplit(prefx, commTag)
	lenRes := len(res)
	if lenRes >= 1 {
		comm = bytes.Replace(res[0], prefx, []byte{}, 1)
	}
	if lenRes > 1 {
		par = res[1:]
	}
	return
}
func (sc *SendContext) DefaultCommandSplit() ([]byte, [][]byte) {
	if !sc.doneDefaultByte {
		sc.commByte, sc.argByte = sc.CommandSplit(CommandPrefx, SpliteTag)
		sc.doneDefaultByte = true
	}
	return sc.commByte, sc.argByte
}
func (sc *SendContext) CommandStringSplit(prefx []byte, commTag []byte) (comm string, par string) {
	c, p := sc.CommandSplit(prefx, commTag)
	if c != nil {
		comm = string(c)
	}
	if p != nil {
		par = string(bytes.Join(p, commTag))
	}
	return
}
func (sc *SendContext) DefaultCommandStringSplit() (comm string, par string) {
	if !sc.doneDefault {
		sc.comm, sc.arg = sc.CommandStringSplit(CommandPrefx, SpliteTag)
		sc.doneDefault = true
	}
	return sc.comm, sc.arg
}
func (sc *SendContext) DefaultHandleCommand(f func(comm string, args []string) bool) {
	sc.HandleCommand(CommandPrefx, SpliteTag, f)
}
func (sc *SendContext) HandleCommand(prefx []byte, commTag []byte, f func(comm string, args []string) bool) {
	c, a := sc.CommandSplit(prefx, commTag)
	comm := string(c)
	argsStrs := make([]string, len(a))
	for i := 0; i < len(a); i++ {
		argsStrs[i] = string(a[i])
	}
	if f(comm, argsStrs) {
		sc.Done()
		sc.SetHandle(true)
	}

}

func (sc *SendContext) Line() SendLine {
	return sc.line
}
func (sc *SendContext) Model() *chat.SendRequest {
	return sc.model
}
func (sc *SendContext) PkgType(ptype int32) {
	sc.model.Pkg.Type = ptype
}
func (sc *SendContext) Data(data []byte) {
	sc.model.Pkg.Data = data
}
func (sc *SendContext) Done() {
	sc.done = true
}
func (sc *SendContext) GetError() error {
	return sc.err
}

func (sc *SendContext) Error(err error) {
	sc.err = err
	sc.Done()
}

type RecvContext struct {
	response *chat.RecvResponse
	err      error
	Appear   bool
	done     bool
}

func (rc *RecvContext) Response() *chat.RecvResponse {
	return rc.response
}
func (rc *RecvContext) EqualType(rtype int32) bool {
	return rc.response.Pkg.Type == rtype
}
func (rc *RecvContext) Done() {
	rc.done = true
}
func (rc *RecvContext) Error(err error) {
	rc.err = err
	rc.Done()
}
func (rc *RecvContext) GetError() error {
	return rc.err
}

type PipManager struct {
	recvPip     []RecvPipHandle
	sendPip     []SendPipHandle
	readyHandle ReadyHandle
}

func NewPipManager() *PipManager {
	pm := new(PipManager)
	pm.recvPip = make([]RecvPipHandle, 0)
	pm.sendPip = make([]SendPipHandle, 0)
	return pm
}
func (pm *PipManager) GetReadyHandle() ReadyHandle {
	return pm.readyHandle
}
func (pm *PipManager) SetReadyHandle(rh ReadyHandle) {
	pm.readyHandle = rh
}
func (pm *PipManager) UseRecv(rp ...RecvPipHandle) *PipManager {
	for _, s := range rp {
		pm.recvPip = append(pm.recvPip, s)
	}
	return pm
}
func (pm *PipManager) UseSend(sp ...SendPipHandle) *PipManager {
	for _, s := range sp {
		pm.sendPip = append(pm.sendPip, s)
	}
	return pm
}
func (pm *PipManager) UseDefaultRecv() *PipManager {
	return pm.UseRecv(func(rc *RecvContext) {
		if rc.Appear {
			fmt.Println(rc.Response().From + " -> " + string(rc.Response().Pkg.Data))
		}
		rc.Done()
	})
}
func (pm *PipManager) UseSendTo(prefx []byte) *PipManager {
	return pm.UseSend(func(rc *SendContext) {
		com, args := rc.DefaultCommandSplit()
		if bytes.Equal(com, prefx) {
			if len(args) > 1 {
				rc.model.To = string(args[0])
				rc.Data(bytes.Join(args[1:], SpliteTag))
				rc.PkgType(1)
				rc.Done()
			} else {
				rc.Error(errors.New("错误参数数目"))
				rc.SetHandle(true)
				rc.Done()
			}
		}
	})
}
func (pm *PipManager) HandleRecv(rep *chat.RecvResponse, appear bool) *RecvContext {
	ctx := &RecvContext{Appear: appear}
	ctx.response = &chat.RecvResponse{Pkg: &chat.SendPkg{}}
	ctx.response = rep
	lenPip := len(pm.recvPip)
	for i := 0; i < lenPip; i++ {
		pm.recvPip[i](ctx)
		if ctx.done {
			break
		}
	}
	return ctx
}
func (pm *PipManager) HandleSend(line SendLine) *SendContext {
	ctx := &SendContext{}
	ctx.model = &chat.SendRequest{Pkg: &chat.SendPkg{Medata: make(map[string][]byte, 0)}}
	ctx.line = line
	lenPip := len(pm.sendPip)
	for i := 0; i < lenPip; i++ {
		pm.sendPip[i](ctx)
		if ctx.done {
			break
		}
	}
	return ctx
}
