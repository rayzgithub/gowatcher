# gowatcher

# 依赖包

go get github.com/howeyc/fsnotify

go get github.com/gomodule/redigo/redis

go get github.com/gpmgo/gopm/modules/log

# 项目说明

监听文件夹及子文件夹下的文件新增，并将新增文件写入至redis

在config/app.go完成配置，配置监听目录，忽略目录、及redis配置信息

可监听修改删除等事件，修改程序中相关代码即可
