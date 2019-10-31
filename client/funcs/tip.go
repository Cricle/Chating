package funcs

import (
	"Chating/client/initer"
	"Chating/client/tiping"
)

const (
	COMMAND_HELP = "help"
)

func TipSend(rc *initer.SendContext) {
	comm, args := rc.DefaultCommandStringSplit()
	if comm == COMMAND_HELP {
		if args == "" {
			tiping.DefaultTiper.RunAll()
		} else {
			tiping.DefaultTiper.Run(args)
		}

		rc.Done()
		rc.SetHandle(true)
	}
}
