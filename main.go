package main

import "fmt"

import "crash_server_win/analyze"
import "crash_server_win/common"
import "crash_server_win/receiver"

func taskCb(info interface{}, succ bool, result string) {
	fmt.Println(succ, result);
}

func main() {
	common.InitLogger();
	analyze.InitAnalyze(taskCb);

	receiver.RunReceiver();
}