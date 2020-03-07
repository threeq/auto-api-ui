package main

import "flag"

type args struct {
	Addr string
	UIPath   string
}
var Args = new(args)

func init() {
	flag.StringVar(&Args.Addr, "addr", ":8080", "服务端地址")
	flag.StringVar(&Args.UIPath, "ui", "asserts/ui", "ui 路径")

}

func CmdArgsParse() {
	flag.Parse()
}
