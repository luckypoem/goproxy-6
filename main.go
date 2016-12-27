package main

import (
	"log"
	"github.com/yuya008/goproxy/server"
	"github.com/mkideal/cli"
	"os"
)

type Args struct {
	cli.Helper
	Mode		string `cli:"mode" usage:"模式local或者remote"`
	Host		string `cli:"host" usage:"本地绑定地址端口"`
	RemoteHost	string `cli:"remote" usage:"远程服务地址端口,本地服务需要配置"`
	AESKey		string `cli:"aeskeypath" usage:"加密密钥文件"`
}

func printArgs(args *Args)  {
	log.Println("本地服务监听主机及端口配置:", args.Host)
	if args.RemoteHost != "" {
		log.Println("远程服务监听主机及端口配置:", args.RemoteHost)
	}
	log.Println("AES加密Key文件路径:%s", args.AESKey)
	if args.Mode == "local" {
		log.Println("启动初始化本地服务器")
	} else {
		log.Println("初始化并启动远程服务")
	}
}

// 入口文件
func main() {
	cli.Run(new(Args), func(ctx *cli.Context) error {
		args := ctx.Argv().(*Args)
		switch args.Mode {
		case "local":
			printArgs(args)
			server.InitLocalServer(&server.LocalService{
				Host : args.Host,
				RemoteHost : args.RemoteHost,
				AESKeyPath : args.AESKey,
			})
		case "remote":
			printArgs(args)
			server.InitRemoteServer(&server.RemoteService{
				Host : args.Host,
				AESKeyPath : args.AESKey,
			})
		default:
			log.Println("参数错误")
			os.Exit(1)
		}
		return nil
	}, "goproxy 翻墙利器")
}
