package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

// 数据库引擎全局变量
var Engine *xorm.Engine



// 连接数据库并同步表结构
func ConnectDB(dsn string) (*xorm.Engine, error) {
	var err error
	Engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := Engine.Ping(); err != nil {
		return nil, err
	}
	if err := Engine.Sync2(&User{}, &Link{}); err != nil {
		return nil, err
	}
	return Engine, nil
}

// 初始化数据库连接
func init() {
	_, err := ConnectDB("root:2470@tcp(81.70.152.142:9000)/shortener_db")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected")
}
