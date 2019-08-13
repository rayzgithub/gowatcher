package main

import (
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"rayz/filewatcher/config"
	"rayz/filewatcher/redis"
)

const sys_all_events = 0xfff

func getDirs(dirPath string) ([]string, error) {
	infos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var dirs []string

	pathSep := string(os.PathSeparator)

	for _, f := range infos {
		if f.IsDir() {
			dirs = append(dirs, dirPath+pathSep+f.Name())
			subDirs, _ := getDirs(dirPath + pathSep + f.Name())
			for _, dir := range subDirs {
				dirs = append(dirs, dir)
			}
		}
	}
	return dirs, nil
}

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	redisclient := redis.Get()

	var ignore = config.Ignores

	fileChan := make(chan string)

	go func() {
		for {
			e := <-fileChan
			//fmt.Println(e)
			_, err := redisclient.Do("LPUSH", "filelist", e)
			if err != nil {
				//panic(err)
				log.Printf(err.Error())
			}
			//fmt.Println(reply)
		}
	}()

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				//文件被创建
				if ev.IsCreate() {
					evname := ev.Name
					//log.Println("name : ", evname)
					if f, _ := os.Stat(evname); f.IsDir() {

						_, fileName := filepath.Split(evname)
						ignoreFlag := false
						for _, v := range ignore {
							if v == fileName {
								ignoreFlag = true
							}
						}

						if !ignoreFlag {
							log.Println("dir create")

							err := watcher.Watch(evname)
							if err != nil {
								log.Println(evname + "开启监控失败")
							} else {
								//log.Println(evname + "开启监控")
							}
						} else {
							log.Println(evname + "已被忽略")
						}

					} else {
						log.Println(evname)
						go func() {
							fileChan <- evname
						}()
					}
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	//获取所有文件夹
	subDirs, _ := getDirs(config.Dir)

	var dirs []string

	dirs = append(dirs, config.Dir)

	dirs = append(dirs, subDirs...)

	for _, dir := range dirs {
		err := watcher.Watch(dir)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(dir + "目录监控已启动")
		}
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}
