package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/gpmgo/gopm/modules/log"
	"rayz/gowatcher/config"
)

var Conn redigo.Conn

func init() {
	c, err := redigo.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		redigo.DialPassword(config.RedisAuth),
		redigo.DialDatabase(config.RedisDbIndex),
	)
	if err != nil{
		log.Error("redis 连接错误")
		panic(err.Error())
	}
	Conn = c
}

func Get() redigo.Conn {
	return Conn
}
