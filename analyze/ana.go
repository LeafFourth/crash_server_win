package analyze


import "crash_server_win/defines"
import "fmt"
import "os/exec"

type AnaCb func(info interface{}, succ bool, result string)
//type AnaCb func(succ bool)

type Task struct {
	Ver string
	File string
	Info interface{}
}

var taskQue chan Task;
var analyzeCb AnaCb;

func getPdbPath(ver string) string {
	priPdbPath := defines.PdbPath + ver;

	pdbPath := `srv*d:\symbolslocal*http://msdl.microsoft.com/download/symbols;` + 
		defines.CommonPdbPath + ";" + 
		priPdbPath;
	return pdbPath;
}

func worker() {
	for {
		t := <- taskQue;

		ana(t.File, t.Ver, analyzeCb, t.Info);
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

func ana(crashFile string, ver string, cb AnaCb, eInfo interface{}) {
	arg0 := "-y"
	arg1 := getPdbPath(ver);
	arg2 := `-z`
	arg3 := crashFile
	arg4 := `-c`
	arg5 := `.ecxr;kb L1000;q`

	cmd := exec.Command(defines.CdbPath, arg0, arg1, arg2, arg3, arg4, arg5);

	out, err1 := cmd.Output();
	if err1 != nil {
		fmt.Println(string(out));
	}

	if out == nil {
		cb(eInfo, false, "");
		return;
	}

	st := extractCallback(out);
	if st == nil {
		cb(eInfo, false, "");
		return;
	}
	cb(eInfo, true, string(st));
}

func findByteInSlice(longData []byte, s byte, index int) int {
	if index < 0 {
		return -1;
	}

	length := len(longData);
	if index >= length {
		return -1;
	}
	for i := index; i < length; i++ {
		if longData[i] == s {
			return i;
		}
	}

	return -1;
}

func findBytesInSlice(longData []byte, s []byte, index int) int {
	if index < 0 {
		return -1;
	}

	length  := len(longData);
	length2 := len(s);

	if index >= length {
		return -1;
	}

	if length2 <= 0 {
		return -1;
	}

	if length - index < length2 {
		return -1;
	}

	var end = length - length2;
	for i := index; i <= end; i++ {
		for j, k := i, 0; ; {
			if k == length2 {
				return i;
			}

			if longData[j] != s[k] {
				break;
			}

			j++;
			k++;
		}
	}
	return -1;
}

func extractCallback(verbose []byte) []byte {
	if verbose == nil || len(verbose) <= 0 {
		return nil;
	}

	beginTag := "ChildEBP";
	endTag   := "quit";
	endPos   := len(verbose) - 1;
	
	validBegin := -1;
	validEnd   := -1;

	var pos int;
	for  {
		pos = findBytesInSlice(verbose, []byte(beginTag), pos);
		if (pos == -1) {
			break;
		}

		pos2 := pos + len(beginTag);
		if pos > 0 && verbose[pos - 1] != '\n' { // beginTag must be begin of a line;
			pos = pos2;
			continue;
		}

		pos = findByteInSlice(verbose, '\n', pos2);
		if pos == -1 || pos >= endPos {
			break;
		}

		validBegin = pos + 1;
		break;
	}

	if validBegin < 0 {
		return nil;
	}

	for {
		pos = findBytesInSlice(verbose, []byte(endTag), pos);
		if (pos == -1) {
			break;
		}

		if pos > 0 && verbose[pos - 1] != '\n' { // beginTag must be begin of a line;
			pos = pos + len(beginTag);
			continue;
		}
		validEnd = pos;
		break;
	}

	if validEnd < 0 {
		validEnd = endPos + 1;
	}

	return verbose[validBegin: validEnd];

}
