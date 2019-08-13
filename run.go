package main

import (
	"log"
	"rayz/filewatcher/listener"
	"rayz/filewatcher/redis"
)

func main() {

	//初始化redis连接
	redisclient := redis.Get()

	//文件新增响应事件
	listener.Instance.OnFileCreate = func(file string) {
		_, err := redisclient.Do("LPUSH", "filelist", file)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	done := make(chan bool)

	listener.Instance.Listen()

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	listener.Instance.Close()
}
