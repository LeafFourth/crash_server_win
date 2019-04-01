package common

import "log"

import "utilities"

var ErrorLogger *log.Logger
var InfoLogger  *log.Logger
var DebugLogger *log.Logger

func InitLogger() {
	writer := utilities.NewLogWriter();
	writer.SetLogPath("E:/log.txt");

	ErrorLogger = log.New(writer, "-E- ", log.Llongfile | log.Lmicroseconds);
	InfoLogger  = log.New(writer, "-I- ", log.Llongfile | log.Lmicroseconds);
	DebugLogger = log.New(writer, "-D- ", log.Llongfile | log.Lmicroseconds);
}