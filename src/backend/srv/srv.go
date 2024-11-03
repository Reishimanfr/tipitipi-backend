package srv

import (
	ovh "bash06/strona-fundacja/src/backend/aws"
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/middleware"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	Db     *gorm.DB
	Log    *zap.Logger
	Ovh    *ovh.Worker
	Router *gin.Engine
	Http   *http.Server
	Argon  *core.Argon2idHash
}

type ServerConfig struct {
	Port       string
	Testing    bool
	HttpConfig *http.Server
	AwsConfig  *aws.Config
	CorsConfig *cors.Config
}

func New(c *ServerConfig) (*Server, error) {
	s := &Server{}

	log, err := core.InitLogger()
	if err != nil {
		return nil, err
	}

	s.Log = log

	db, err := core.InitDb(c.Testing)
	if err != nil {
		return nil, err
	}

	s.Db = db

	if c.AwsConfig == nil {
		return nil, fmt.Errorf("aws config is nil")
	}

	ses, err := session.NewSession(c.AwsConfig)
	if err != nil {
		return nil, err
	}

	s3Cli := s3.New(ses)

	s.Ovh = &ovh.Worker{
		Session: ses,
		S3:      s3Cli,
	}

	s.Router = gin.New(func(e *gin.Engine) {
		e.Use(middleware.RateLimiterMiddleware(middleware.NewRateLimiter(5, 10)))

		// Only set CORS headers if we're not testing
		if !c.Testing {
			e.Use(cors.New(*c.CorsConfig))
		}
	})

	c.HttpConfig.Addr = ":" + c.Port
	c.HttpConfig.Handler = s.Router

	s.Http = c.HttpConfig

	s.Argon = core.NewArgon2idHash(1, 32, 64*1024, 32, 256)

	return s, nil
}
