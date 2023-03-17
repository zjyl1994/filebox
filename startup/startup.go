package startup

import (
	_ "embed"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/zjyl1994/filebox/server"
	"github.com/zjyl1994/filebox/vars"
)

//go:embed brand.txt
var brandLogo string

func Main() error {
	flag.StringVar(&vars.DataDir, "dir", ".", "数据目录")
	flag.StringVar(&vars.Listen, "listen", "localhost:9496", "监听地址")
	flag.StringVar(&vars.Username, "username", "admin", "WebDAV 用户名")
	flag.StringVar(&vars.Password, "password", "admin", "WebDAV 密码")
	flag.StringVar(&vars.Title, "title", "文件盒子", "网站标题")
	flag.StringVar(&vars.CorsOrigin, "cors", "", "CORS 允许的来源")
	flag.Parse()
	// init vars
	dir, err := filepath.Abs(vars.DataDir)
	if err != nil {
		return err
	}
	vars.DataDir = dir
	// print brand
	fmt.Println(brandLogo)
	fmt.Println("Listen on", vars.Listen)
	fmt.Println("Data dir", vars.DataDir)
	return server.Run()
}
