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

	task := analyze.Task{Ver: "1.0", File: `E:\TestC.dmp`};
	analyze.RunTask(task);

	receiver.InitReceiver();
}