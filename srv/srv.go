package srv

import (
	"bash06/tipitipi-backend/core"
	"bash06/tipitipi-backend/flags"
	"bash06/tipitipi-backend/middleware"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	Db     *gorm.DB
	Log    *zap.Logger
	Router *gin.Engine
	Http   *http.Server
	Argon  *core.Argon2idHash
}

type ServerConfig struct {
	CorsConfig *cors.Config
	HttpConfig *http.Server
}

type FileUploadResult struct {
	Mimetype string
	Size     int64
	Filename string
}

func New(c *ServerConfig) (*Server, error) {
	log, err := core.InitLogger()
	if err != nil {
		return nil, err
	}
	db, err := core.InitDb()
	if err != nil {
		return nil, err
	}

	s := &Server{
		Log:    log,
		Db:     db,
		Router: gin.Default(),
		Argon:  core.NewArgon2idHash(1, 32, 64*1024, 32, 256),
	}

	s.Router.Use(middleware.RateLimiterMiddleware(middleware.NewRateLimiter(5, 10)))
	s.Router.Use(cors.New(*c.CorsConfig))

	c.HttpConfig = &http.Server{
		Addr:    ":" + *flags.Port,
		Handler: s.Router.Handler(),
	}

	return s, nil
}

func (s *Server) DownloadFile(f *multipart.FileHeader) (*FileUploadResult, error) {
	dirPath := filepath.Join(flags.BasePath, "files")

	// Check if dir exists and create if it doesn't
	if _, err := os.Stat(dirPath); errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create files directory at %s: %v", dirPath, err)
		}
	}

	srcFile, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open source file %s: %v", f.Filename, err)
	}

	defer srcFile.Close()

	b := make([]byte, 512)
	if _, err := srcFile.Read(b); err != nil {
		return nil, fmt.Errorf("failed to read file header for %s: %v", f.Filename, err)
	}

	srcFile.Seek(0, 0)

	result := &FileUploadResult{
		Mimetype: mimetype.Detect(b).String(),
		Size:     f.Size,
		Filename: core.RandStr(10) + filepath.Ext(f.Filename),
	}

	dest, err := os.Create(filepath.Join(dirPath, result.Filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file for %s: %v", f.Filename, err)
	}

	defer dest.Close()

	writer := bufio.NewWriter(dest)
	if _, err := io.Copy(writer, srcFile); err != nil {
		return nil, fmt.Errorf("failed to copy file contents to destination for %s: %v", f.Filename, err)
	}

	defer writer.Flush()

	return result, nil
}

func (s *Server) DeleteFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("no filename provided")
	}

	filePath := filepath.Join(flags.BasePath, "files", filename)

	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("file %s doesn't exist", filename)
		}

		if errors.Is(err, fs.ErrPermission) {
			return fmt.Errorf("insufficient permissions to delete file %s", filename)
		}

		return fmt.Errorf("something went wrong while deleting file %s: %v", filename, err)
	}

	return nil
}
