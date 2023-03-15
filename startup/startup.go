package startup

import (
	"flag"

	"github.com/zjyl1994/filebox/server"
	"github.com/zjyl1994/filebox/vars"
)

func Main() error {
	flag.StringVar(&vars.DataDir, "dir", ".", "数据目录")
	flag.StringVar(&vars.Listen, "listen", ":9496", "监听地址")
	flag.StringVar(&vars.Username, "username", "admin", "WebDAV 用户名")
	flag.StringVar(&vars.Password, "password", "admin", "WebDAV 密码")
	flag.StringVar(&vars.Title, "title", "文件盒子", "网站标题")
	flag.Parse()
	return server.Run()
}
