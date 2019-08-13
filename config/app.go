package config

var (
	//监听的目录
	Dir = `F:\test`
	//Dir     = `/opt/ftp/out`
	Ignores = []string{".git", ".idea"}
	//redis配置
	RedisHost    = "127.0.0.1"
	RedisPort    = 6379
	RedisAuth    = "abcdefg"
	RedisDbIndex = 15
)
