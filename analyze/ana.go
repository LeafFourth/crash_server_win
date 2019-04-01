package analyze


import "crash_server_win/defines"
import "fmt"
import "os/exec"
import "unsafe"

type AnaCb func(info interface{}, succ bool, result string)
//type AnaCb func(succ bool)

type Task struct {
	ver string
	file string
	info interface{}
}

var taskQue chan Task;
var analyzeCb AnaCb;

func getPdbPath(ver string) string {
	return defines.PdbPath + ver;
}

func worker() {
	for {
		t := <- taskQue;

		Ana(t.file, t.ver, analyzeCb, t.info);
	}
}

func InitAnalyze(cb AnaCb) {
	taskQue   = make(chan Task, 8);
	analyzeCb = cb;

	for i := 0; i < 4; i++ {
		go worker();
	}
}

func RunTask(t Task) {
	taskQue <- t;
}

func Ana(crashFile string, ver string, cb AnaCb, eInfo interface{}) {
	priPdbPath := getPdbPath(ver)
	arg0 := "-y"
	arg1 := `srv*d:\symbolslocal*http://msdl.microsoft.com/download/symbols;` + priPdbPath
	arg2 := `-z`
	arg3 := crashFile
	arg4 := `-v`
	arg5 := `-c`
	arg6 := `".ecxr;kb L1000;q"`

	var p unsafe.Pointer;
	fmt.Println(p);

	cmd := exec.Command(defines.CdbPath, arg0, arg1, arg2, arg3, arg4, arg5, arg6);

	out, err1 := cmd.Output();
	if err1 != nil {
		
		cb(eInfo, false, "")
		return;
	}

	fmt.Println(out);
}