package db

import (
	"context"
	"fastlink/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

// 数据库引擎全局变量
var SQLEngine *xorm.Engine
var RedisClient *redis.Client
var Ctx = context.Background()

// 连接数据库并同步表结构
func ConnectMysql() error {
	var err error
	SQLEngine, err = xorm.NewEngine("mysql", config.Global.Mysql.Dsn)
	if err != nil {
		return err
	}
	if err := SQLEngine.Ping(); err != nil {
		return err
	}

	if err := SQLEngine.Sync2(&User{}, &Link{}); err != nil {
		return err
	}
	return nil
}

func ConnectRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Global.Redis.Addr,
		Password: config.Global.Redis.Password,
		DB:       config.Global.Redis.DB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	return err
}

// 初始化数据库连接
func init() {
	err := ConnectMysql()
	if err != nil {
		panic(err)
	}

	err = ConnectRedis()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected")

}
