package flags

import (
	"flag"
	"os"
	"path/filepath"
)

var (
	Dev          = flag.Bool("dev", true, "Enables development logging")
	Port         = flag.String("port", "8080", "Port on which the server should run")
	Secure       = flag.Bool("secure", false, "Enables https")
	CertFilePath = flag.String("cert-file-path", "", "Path to the SSL Cert file")
	KeyFilePath  = flag.String("key-file-path", "", "Path to the SSL key file")

	execPath, _ = os.Executable()
	BasePath    = filepath.Dir(execPath)
)
