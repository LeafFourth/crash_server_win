package defines

//import "fmt"

import "encoding/json"
import "fmt"
import "os"

const confName = "conf.json";

type ports struct {
	ReceiveHttp uint;
}

type paths struct {
	CdbPath        string;
	PdbRoot        string;
	CommmonPdbRoot string;

	TmpDmpZipRoot  string;
	DmpRoot        string;

	ResRoot        string;
	LogsRoot       string;
}

type names struct {
	DmpName     string;
	DmpDescName string;
}

type conf struct {
	Ports ports;
	Paths paths;
	Names names;
}

func createPaths() {
	if err := os.MkdirAll(PdbPath, 0644); err != nil {
		fmt.Println(err);
	}

	if err := os.MkdirAll(CommonPdbPath, 0644); err != nil {
		fmt.Println(err);
	}

	if err := os.MkdirAll(LocalStorePath, 0644); err != nil {
		fmt.Println(err);
	}

	if err := os.MkdirAll(UnzipPath, 0644); err != nil {
		fmt.Println(err);
	}

	if err := os.MkdirAll(LogsRoot, 0644); err != nil {
		fmt.Println(err);
	}
}

func InitDefines(confPath string) {
	f, err := os.Open(confPath);
	if err != nil {
		fmt.Println(err);
		return;
	}

	d := json.NewDecoder(f);
	var c conf;
	d.Decode(&c);

	CdbPath        = c.Paths.CdbPath;
	PdbPath        = c.Paths.PdbRoot;
	CommonPdbPath  = c.Paths.CommmonPdbRoot;
	ResRoot        = c.Paths.ResRoot;
	LocalStorePath = c.Paths.TmpDmpZipRoot;
	UnzipPath      = c.Paths.DmpRoot;
	LogsRoot       = c.Paths.LogsRoot;

	ReceiverPort = c.Ports.ReceiveHttp;

	DmpName     = c.Names.DmpName;
	DmpDescName = c.Names.DmpDescName;

	createPaths();

}