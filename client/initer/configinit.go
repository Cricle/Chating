package initer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	CONFIG_FILE_NAME = "config.cfg"
)

type Initer struct {
	inited    bool
	hasConfig bool
	config    *Config
}

func NewIniter() *Initer {
	i := Initer{
		config: &Config{},
	}
	return &i
}
func (i *Initer) Config() *Config {
	return i.config
}
func (i *Initer) IsInited() bool {
	return i.inited
}
func (i *Initer) Init(stdin *bufio.Reader) (err error) {
	if !i.inited {
		fs, _ := os.Stat(CONFIG_FILE_NAME)
		if fs != nil {
			fi, err := os.Open(CONFIG_FILE_NAME)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			defer fi.Close()
			bs, err := ioutil.ReadAll(fi)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			err = json.Unmarshal(bs, &i.config)

			i.hasConfig = true
			i.inited = true
		} else {
			fmt.Print("IP:")
			line, _, err := stdin.ReadLine()
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			i.config.Address = string(line)
			bs, err := json.Marshal(i.config)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			err = ioutil.WriteFile(CONFIG_FILE_NAME, bs, os.ModePerm)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
		}
	}
	return nil

}
func (i *Initer) WithConfig(stdin *bufio.Reader) {
	fmt.Print("Name is :")
	line, _, err := stdin.ReadLine()
	i.config.Name = string(line)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Print("Pwd is :")
	line, _, err = stdin.ReadLine()
	i.config.Pwd = string(line)
	if err != nil {
		fmt.Println(err.Error())
	}
}
