package server

import "bytes"
import "encoding/json"
import "fmt"
import "io/ioutil"
import "mime/multipart"
import "net/http"
import "os"
import "path/filepath"
import "strconv"
import "strings"

import "utilities"

import "crash_server_win/analyze"
import "crash_server_win/common"
import "crash_server_win/defines"

type dmpDesc struct {
	Ver  string;
	Date string;
	Uid  int;
	Name string;
}

type unzipTask struct {
	zipFile string;
	dst     string;
	cInfo interface {};
	cb    func(unzipTask);
}

var handler *utilities.RequestHandler;

var unzipQue chan unzipTask;
var gateClient *http.Client;

func taskCb(info interface{}, succ bool, result string) {
	bodyBuffer := new(bytes.Buffer);
	bodyWriter := multipart.NewWriter(bodyBuffer);

	w, err := bodyWriter.CreateFormField(defines.CallstackKey);
	if err != nil {
		common.ErrorLogger.Print(info, err);
		return;
	}
	if w == nil {
		common.ErrorLogger.Print(info, "writer nil");
		return;
	}
	w.Write([]byte(result));

	w2, err2 := bodyWriter.CreateFormField(defines.EInfoKey);
	if err != nil {
		common.ErrorLogger.Print(info, err2);
		return;
	}
	if w2 == nil {
		common.ErrorLogger.Print(info, "writer nil");
		return;
	}
	e := json.NewEncoder(w2);
	err3 := e.Encode(info);
	if err3 != nil {
		common.ErrorLogger.Print(info, "writer error:", err3);
		return;
	}

	bodyWriter.Close();

	rq, err4 := http.NewRequest(http.MethodPost, defines.GateSvr + defines.CallstackApi, bodyBuffer);
	if err4 != nil {
		common.ErrorLogger.Print(info, err3);
		return;
	}
	rq.Header.Add("Content-Type", "multipart/form-data; " + "boundary=" + bodyWriter.Boundary());

	{
		reps, err := gateClient.Do(rq);
		//fmt.Println(reps, err);
		if err != nil {
			common.ErrorLogger.Print(info, err);
			return;
		}

		if reps.StatusCode != http.StatusOK {
			common.ErrorLogger.Print(info, "response code:", reps.StatusCode);
			return;
		}
	}
	
}

func handleDefaultPage(w http.ResponseWriter, r *http.Request) bool {
  if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.html";
		http.DefaultServeMux.ServeHTTP(w, r);
		return true;
	}

	return false;
}

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("require ", r.URL.Path);
	if handleDefaultPage(w, r) {
		return;
	}

	path := filepath.Join(defines.ResRoot, r.URL.Path[1:]);
	f, err := os.Open(path);
	if err != nil {
		fmt.Println(err);
		w.WriteHeader(404);
		w.Write([]byte(""));
		return;
	}

	data, err2 := ioutil.ReadAll(f);
	if err2 != nil {
		fmt.Println("read err");
		fmt.Println(err2);
		w.WriteHeader(404);
		w.Write([]byte(""));
		return;
	}
	
	w.Write(data);
}

func receiveCrashFile(w http.ResponseWriter, r *http.Request) {
	osFile, fh, err := r.FormFile("crashFile");
	if err != nil {
		common.ErrorLogger.Print(err);
		w.WriteHeader(http.StatusBadRequest);
		return;
	}

	if fh == nil {
		common.ErrorLogger.Print("no file");
		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("empty"));
		return;
	}

	if osFile == nil {
		common.ErrorLogger.Print("no file");
		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("file error"));
		return;
	}

	common.InfoLogger.Print("received:", fh.Filename);

	filePath := filepath.Join(defines.LocalStorePath, fh.Filename);
	err2 := utilities.WriteFile(osFile, filePath);
	if err2 != nil {
		common.ErrorLogger.Print(err2);
		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("file error"));
		return;
	}
	osFile.Close();

	postUnzipTask(filePath, defines.UnzipPath);
}

func initHandler() {
	handlers := make(map[string]func(http.ResponseWriter, *http.Request));
	handlers["/"] = defaultHandle;
	handlers["/postCrash"] = receiveCrashFile;
	handlers["/postPdbs"] = receivePdbs;

	handler = utilities.NewRequestHandler(&handlers);
}

func runHttpServer() {
	port := ":" + strconv.FormatUint(uint64(defines.ReceiverPort), 10) ;
	err := http.ListenAndServe(port, handler);
	if err != nil {
		common.ErrorLogger.Print(err);
		return;
	}
}

func unzipWorker() {
	for  {
		t := <- unzipQue;
		if len(t.zipFile) <= 0 || len(t.dst) <= 0 {
			continue;
		}

		err3 := utilities.UnzipFile(t.zipFile, t.dst);
		if err3 != nil {
			common.ErrorLogger.Print(t.zipFile, ":", err3);
		}

		if err := os.Remove(t.zipFile); err != nil {
			common.ErrorLogger.Print(t.zipFile, ":", err3);
		}

		if t.cb != nil {
			t.cb(t);
		}
	}
	
}

func initUnzipTask() {
	unzipQue = make(chan unzipTask, 1);

	go unzipWorker();
}

func onUnzipDown(t unzipTask) {
	dirName := strings.TrimSuffix(filepath.Base(t.zipFile), ".zip");
	descName := filepath.Join(t.dst, dirName, defines.DmpDescName);
	dmpName  := filepath.Join(t.dst, dirName, defines.DmpName);

	file, err := os.Open(descName);
	if err != nil {
		common.ErrorLogger.Print("open desc file error:", err);
		return;
	}
	defer file.Close();

	decoder := json.NewDecoder(file);
	var desc dmpDesc;
	decoder.Decode(&desc);
	if len(desc.Ver) <= 0 {
		common.ErrorLogger.Print("no ver info");
		return;
	}
	desc.Name = dirName;

	task := analyze.Task{
				Ver: desc.Ver,
				File: dmpName,
				Info: &desc,
			}

	analyze.RunTask(task);
}

func postUnzipTask(zipFile, dst string) {
	task := unzipTask{zipFile: zipFile,
					  dst: dst,
					  cb: onUnzipDown,
					  cInfo: nil};
	unzipQue <- task;
}

func RunReceiver () {
	analyze.InitAnalyze(taskCb);
	initUnzipTask();

	gateClient = new(http.Client);

	initHandler();
	runHttpServer();
}

func receivePdbs(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1000000);
	numV := r.MultipartForm.Value["pdb_num"];
	if numV == nil || len(numV) <= 0 {
		common.ErrorLogger.Print("pdb num error");

		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("pdb num error"));
		return;
	}

	num, err := strconv.ParseInt(numV[0], 10, 0);
	if err != nil {
		common.ErrorLogger.Print("num format error:", numV[0]);

		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("num format error"));
		return;
	}


	verV := r.MultipartForm.Value["ver"];
	if verV == nil || len(verV) <= 0 {
		common.ErrorLogger.Print("ver error");

		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("ver error"));
		return;
	}
	root := defines.PdbPath + verV[0] + `\`;

	err2 :=  os.MkdirAll(root, 0644);
	if err2 != nil {
		common.ErrorLogger.Print(err2);

		w.WriteHeader(http.StatusInternalServerError);
		w.Write([]byte("unknown error"));
		return;
	}

	for i := int64(0); i < num; i++ {
		fileKey := "pdb" + strconv.FormatInt(i, 10);

		fhV := r.MultipartForm.File[fileKey];
		if fhV == nil ||  len(fhV) <= 0 {
			common.ErrorLogger.Print("file not exist:", fileKey);
			continue;
		}

		file := fhV[0];
		if file == nil {
			common.ErrorLogger.Print("file not exist:", fileKey);
			continue;
		}

		osFile, err3 := file.Open();
		if err3 != nil {
			common.ErrorLogger.Print("file open error:", fileKey);
			continue;
		}
		defer osFile.Close();


		filePath := root + file.Filename;
		if err4 := utilities.WriteFile(osFile, filePath); err4 != nil {
			common.ErrorLogger.Print("file save error:", fileKey);
			continue;
		}
	}
}