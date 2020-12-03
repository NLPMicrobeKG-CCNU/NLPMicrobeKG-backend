package model

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/log"
)

type Database struct {
	Redis *redis.Pool
}

var DB *Database

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Error("Open database failed",
			zap.String("reason", err.Error()),
			zap.String("detail", fmt.Sprintf("Database connection failed. Database name: %s", name)))
	}

	// set for db connection
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

func openRedis(server, password string) (conn *redis.Pool) {
	pool := &redis.Pool{
		// 初始化链接数量
		MaxIdle:     200,
		MaxActive:   0,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return DialRedis(server, password)
		},
	}

	return pool
}

func InitRedis() *redis.Pool {
	server := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	conn := openRedis(server, password)
	return conn
}

func GetRedis() *redis.Pool {
	return InitRedis()
}

func DialRedis(server, password string) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", server)
	if err != nil {
		log.Error("Open redis failed",
			zap.String("reason", err.Error()),
			zap.String("detail", fmt.Sprintf("Redis connection failed. Database name: %s", server)))
		return nil, err
	}

	return conn, nil
}

func (db *Database) Init() {
	DB = &Database{
		Redis: GetRedis(),
	}
}

func (db *Database) RedisClose() {
	DB.Redis.Close()
}

func (db *Database) Close() {
	db.RedisClose()
}
