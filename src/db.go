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
	ID         int       `xorm:"'link_id'         pk autoincr"`
	SourceUrl  string    `xorm:"'source_url'      notnull"`
	ShortUrl   string    `xorm:"'short_url'       index notnull unique"`
	UserID     int       `xorm:"'user_id'         notnull"`
	ClickCount int       `xorm:"'click_count'     default(0)"`
	CreatedAt  time.Time `xorm:"'created_at'      created"`
	ExpireAt   time.Time `xorm:"'expire_at'       notnull"`
}

func ConnectDB(dsn string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := engine.Ping(); err != nil {
		return nil, err
	}
	if err := engine.Sync2(&User{}, &Link{}); err != nil {
		return nil, err
	}
	return engine, nil
}


func init() {
	_, err := ConnectDB("root:2470@tcp(81.70.152.142:9000)/shortener_db")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected")
}
