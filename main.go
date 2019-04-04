package main


import "crash_server_win/common"
import "crash_server_win/receiver"

func main() {
	common.InitLogger();

	receiver.RunReceiver();
}