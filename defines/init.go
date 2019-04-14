package defines

//import "fmt"

import "encoding/json"
import "fmt"
import "os"

const confName = "conf.json";

type servers struct {
	ReceiveHttp    uint;

	GateSvr 	   string;
	CallstackApi   string;
	CallstackKey   string;
	EInfoKey       string;
}

type paths struct {
	CdbPath        string;
	PdbRoot        string;
	CommmonPdbRoot string;

	DmpZipRoot     string;
	DmpUnzipRoot   string;

	ResRoot        string;
	LogsRoot       string;
}

type names struct {
	DmpName     string;
	DmpDescName string;
}

type conf struct {
	Servers servers;
	Paths   paths;
	Names   names;
}

func createPaths() {
	if err := os.MkdirAll(PdbPath, 0644); err != nil {
		fmt.Println(PdbPath, ":", err);
	}

	if err := os.MkdirAll(CommonPdbPath, 0644); err != nil {
		fmt.Println(CommonPdbPath, ":", err);
	}

	if err := os.MkdirAll(DmpZipRoot, 0644); err != nil {
		fmt.Println(DmpZipRoot, ":", err);
	}

	if err := os.MkdirAll(DmpUnzipRoot, 0644); err != nil {
		fmt.Println(DmpUnzipRoot, ":", err);
	}

	if err := os.MkdirAll(LogsRoot, 0644); err != nil {
		fmt.Println(LogsRoot, ":", err);
	}
}

func InitDefines(confPath string) {
	f, err := os.Open(confPath);
	if err != nil {
		fmt.Println(err);
		return;
	}
	defer f.Close();

	d := json.NewDecoder(f);
	var c conf;
	if err := d.Decode(&c); err != nil {
		fmt.Println(err);
		return;
	}

	CdbPath        = c.Paths.CdbPath;
	PdbPath        = c.Paths.PdbRoot;
	CommonPdbPath  = c.Paths.CommmonPdbRoot;
	ResRoot        = c.Paths.ResRoot;
	DmpZipRoot     = c.Paths.DmpZipRoot;
	DmpUnzipRoot   = c.Paths.DmpUnzipRoot;
	LogsRoot       = c.Paths.LogsRoot;

	ReceiverPort = c.Servers.ReceiveHttp;
	GateSvr      = c.Servers.GateSvr;
	CallstackApi = c.Servers.CallstackApi;
	CallstackKey = c.Servers.CallstackKey;
	EInfoKey     = c.Servers.EInfoKey;

	DmpName     = c.Names.DmpName;
	DmpDescName = c.Names.DmpDescName;

	createPaths();

}