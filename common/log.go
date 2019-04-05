package common

import "log"
import "path/filepath"
import "utilities"

import "crash_server_win/defines"

var ErrorLogger *log.Logger
var InfoLogger  *log.Logger
var DebugLogger *log.Logger

func InitLogger() {
	writer := utilities.NewLogWriter();
	writer.SetLogPath(filepath.Join(defines.LogsRoot, "log.txt"));

	ErrorLogger = log.New(writer, "-E- ", log.Llongfile | log.Lmicroseconds);
	InfoLogger  = log.New(writer, "-I- ", log.Llongfile | log.Lmicroseconds);
	DebugLogger = log.New(writer, "-D- ", log.Llongfile | log.Lmicroseconds);
}