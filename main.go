package main


import "crash_server_win/common"
import "crash_server_win/defines"
import "crash_server_win/server"

func main() {
	defines.InitDefines("E:/code/go/src/crash_server_win/configure/conf.json");
	common.InitLogger();

	server.RunReceiver();
}