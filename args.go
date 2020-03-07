package main

import "flag"

type args struct {
	Addr          string
	UIPath        string
	Host          string
	Path          string
	Title         string
	Description   string
	Contact       string
	TermOfService string
}

var Args = new(args)

func init() {
	flag.StringVar(&Args.Addr, "addr", ":8080", "服务端地址")
	flag.StringVar(&Args.UIPath, "ui", "asserts/ui", "ui 路径")
	flag.StringVar(&Args.Host, "host", "127.0.0.1:8080", "服务地址")
	flag.StringVar(&Args.Path, "path", "/docs-ui", "服务路径")
	flag.StringVar(&Args.Title, "title", "文档服务", "文档标题")
	flag.StringVar(&Args.Description, "description", "自动化微服务文档服务", "文档描述")
	flag.StringVar(&Args.Contact, "contact", "", "文档中联系我们")
	flag.StringVar(&Args.TermOfService, "term-service", "", "文档团队服务")
}

func CmdArgsParse() {
	flag.Parse()
}
