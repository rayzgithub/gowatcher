package listener

import (
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"rayz/gowatcher/config"
)

type listener struct {
	dir          string
	ignore       []string
	watcher      *fsnotify.Watcher
	fileChan     chan string
	OnFileCreate func(string)
}

var Instance listener

func init() {
	//初始化目录
	Instance.dir = config.Dir
	//初始化忽略目录
	Instance.ignore = config.Ignores
	//初始化监听
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	Instance.watcher = watcher
}

func (l *listener) onDirCreate(s string) {
	//是否是目录
	_, fileName := filepath.Split(s)
	ignoreFlag := inSlice(fileName, l.ignore)

	if !ignoreFlag {
		//log.Println("dir create")
		//设置监听此目录
		err := l.watcher.Watch(s)
		if err != nil {
			log.Println(s + "开启监控失败，message：" + err.Error())
		} else {
			//开启监控成功
		}
	} else {
		log.Println(s + "已被忽略")
	}
}

func (l *listener) Listen() {
	//开启goroutine 接收channel数据
	go func() {
		for {
			e := <-l.fileChan
			l.OnFileCreate(e)
		}
	}()

	// 开启goroutine 监听事件
	go func() {
		for {
			select {
			case ev := <-l.watcher.Event:
				evname := ev.Name
				//文件被创建
				if ev.IsCreate() {
					if f, _ := os.Stat(evname); f.IsDir() {
						//目录被创建
						l.onDirCreate(evname)
					} else {
						//文件被创建
						log.Println(evname)
						go func() {
							l.fileChan <- evname
						}()
					}
				}
			case err := <-l.watcher.Error:
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
		err := l.watcher.Watch(dir)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(dir + "目录监控已启动")
		}
	}
}

func (l *listener) Close() {
	l.watcher.Close()
}

func inSlice(s string, sli []string) bool {
	var flag = false
	for _, v := range sli {
		if v == s {
			flag = true
		}
	}
	return flag
}

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
