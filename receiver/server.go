package receiver

import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strconv"
import "strings"

import "crash_server_win/common"
import "crash_server_win/defines"

func handleDefaultPage(w http.ResponseWriter, r *http.Request) bool {
  if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.html";
		http.DefaultServeMux.ServeHTTP(w, r);
		return true;
	}

	return false;
}

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm();
	// token := r.Form["token"];
	// if token == nil {
	// 	fmt.Println("args error");
	// 	w.WriteHeader(401);
	// 	w.Write([]byte(""));
	// 	return;
	// }

	// if _, e := auth.CheckToken(token[0]); e != nil {
	// 	fmt.Println(e);
	// 	w.WriteHeader(401);
	// 	w.Write([]byte(""));
	// 	return;
	// }

	fmt.Println("require ", r.URL.Path);
	if handleDefaultPage(w, r) {
		return;
	}



  path := defines.ResRoot + r.URL.Path[1:];
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
	r.ParseMultipartForm(1000000);
	files := r.MultipartForm.File["crashFile"];
	if files == nil || len(files) <= 0 {
		fmt.Println(1);
		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("empty"));
		return;
	}
	file := files[0];
	if file == nil {
		fmt.Println(2);
		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("empty"));
		return;
	}

	osFile, err := file.Open();
	defer osFile.Close();
	if err != nil {
		fmt.Println(3, err);
		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("file error"));
		return;
	}

	filePath := defines.LocalStorePath + file.Filename;
	err2 := common.WriteFile(osFile, filePath);
	if err2 != nil {
		fmt.Println(3, err2);
		w.WriteHeader(http.StatusBadRequest);
		w.Write([]byte("file error"));
		return;
	}

	err3 := common.UnzipFile(filePath, "E:/tmp");
	if err3 != nil {
		fmt.Println(err3);
	}
}

func initRouter() {
	http.HandleFunc("/", defaultHandle);
	http.HandleFunc("/postCrash", receiveCrashFile);
}

func initHttpServer() {
	port := ":" + strconv.FormatUint(uint64(defines.ReceiverPort), 10) ;
	err := http.ListenAndServe(port, nil);
	if err != nil {
		common.ErrorLogger.Print(err);
		return;
	}


}

func RunReceiver () {
	initRouter();
	initHttpServer();
}