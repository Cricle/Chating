package tiping

import (
	"fmt"

	color "github.com/gookit/color"
)

var (
	DefaultTiper = NewTiper()
	EmptyDesc    = func() {}
)

type DescriptHandle func()

type Tiper struct {
	items map[string]DescriptHandle
}

func NewTiper() *Tiper {
	t := new(Tiper)
	t.items = make(map[string]DescriptHandle, 0)
	return t
}

func (t *Tiper) Add(name string, handle DescriptHandle) {
	hd := t.items[name]
	if hd != nil {
		panic("重复命令:" + name)
	}
	t.items[name] = handle
}
func (t *Tiper) Get(name string) DescriptHandle {
	return t.items[name]
}
func (t *Tiper) Walk(walkFunc func(string, DescriptHandle) bool) {
	for key, val := range t.items {
		if walkFunc(key, val) {
			break
		}
	}
}
func (t *Tiper) RunAll() {
	t.Walk(func(name string, handle DescriptHandle) bool {
		handle()
		return false
	})
}
func (t *Tiper) Run(name string) {
	t.Walk(func(n string, handle DescriptHandle) bool {
		if name == n {
			handle()
		}
		return false
	})
}

type CmdArg struct {
	Name     string
	Descript string
}

func MakeDefaultPrint(c color.Color, name, desciprt string, args ...CmdArg) func() {
	return func() {
		c.Printf("Command[%s]\n", name)
		color.FgGray.Printf("Descript[%s]\n", desciprt)
		if len(args) != 0 {
			color.FgMagenta.Println("Paramters:")
			for i := 0; i < len(args); i++ {
				arg := args[i]
				fmt.Printf("Name[%s]\n\tDescript[%s]\n", arg.Name, arg.Descript)
			}
			fmt.Println()
		}
		fmt.Println("---------------------------")
	}
}
