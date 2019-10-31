package data

import "time"

type User struct {
	Id         int64     `xorm:"'id' pk autoincr index notnull"`
	Name       string    `xorm:"varchar(32) unique notnull 'name'"`
	Pwd        string    `xorm:"varchar(128) notnull 'pwd'"`
	CreateTime time.Time `xorm:"'create_time' created notnull"`
}
