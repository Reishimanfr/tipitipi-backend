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
	CertFilePath = flag.String("cert-file-path", "", "(absolute) Path to SSL certificate file")
	KeyFilePath  = flag.String("key-file-path", "", "(absolute) Path to SSL key file")
	TokenSize    = flag.Int("token-size", 128, "The amount of bytes to use for tokens")

	execPath, _ = os.Executable()
	BasePath    = filepath.Dir(execPath)
)
