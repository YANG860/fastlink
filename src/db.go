package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

type User struct {
	ID           int       `xorm:"'user_id'     pk  autoincr"`
	RegisteredAt time.Time `xorm:"'created_at'  created"`
	Account      string    `xorm:"'account'     notnull"`
	PwHash       string    `xorm:"'pw_hash'     notnull"`
	Valid        bool      `xorm:"'valid'       default(1) notnull"`

	LinkCount       int       `xorm:"'link_count'      default(0)"`
	LatestCreatedAt time.Time `xorm:"'latest_created_at'"`
}

type Link struct {
	ID        int       `xorm:"'link_id'         pk autoincr"`
	SourceUrl string    `xorm:"'source_url'      notnull"`
	ShortUrl  string    `xorm:"'short_url'       index notnull unique"`
	Userid    int       `xorm:"'user_id'         notnull"`
	CreatedAt time.Time `xorm:"'created_at'      created"`
	ExpireAt  time.Time `xorm:"'expire_at'       notnull"`
}

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:2470@tcp(81.70.152.142:9000)/shortener_db")
	if err != nil {
		panic(err)
	}
	err = engine.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("success!!!")

	engine.Sync2(&User{})
	engine.Sync2(&Link{})
}
