package routers

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/uber/jaeger-client-go"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"myself_framwork/configs"
	"myself_framwork/library/httpserver/ginserver"
	loger "myself_framwork/library/logger"
	"myself_framwork/library/tracing"
	//!! V1 API List
)

// ServerEnv ..
type ServerEnv struct {
	ServiceName     string `env:"API_SERVICE" default:"API_SERVICE"`
	OpenTracingHost string `env:"OPEN_TRACING_HOST" default:"0.0.0.0:4000"`
	DebugMode       string `env:"DEBUG_MODE" default:"debug"`
	ReadTimeout     int    `env:"READ_TIMEOUT" default:"120"`
	WriteTimeout    int    `env:"WRITE_TIMEOUT" default:"120"`
	AppEnv          string `env:"APP_ENV" default:"dev"`
}

var (
	server ServerEnv
)

// Server ..
func Server(listenAddress string) error {
	sugarLogger := loger.GetLogger()

	ThisRouter := ThisRouter{}
	ThisRouter.InitTracing()
	ThisRouter.Routers()
	defer ThisRouter.Close()

	err := ginserver.GinServerUp(listenAddress, ThisRouter.Router)

	if err != nil {
		fmt.Println("Error:", err)
		sugarLogger.Error("Error ", zap.Error(err))
		return err
	}

	fmt.Println("Server UP")
	sugarLogger.Info("Server UP ", zap.String("listenAddress", listenAddress))

	return err
}

// ThisRouter ..
type ThisRouter struct {
	Tracer   opentracing.Tracer
	Reporter jaeger.Reporter
	Closer   io.Closer
	Err      error
	GinFunc  gin.HandlerFunc
	Router   *gin.Engine
}

// Routers ..
func (ThisRouter *ThisRouter) Routers() {
	if server.AppEnv == "DEV" {
		gin.SetMode(server.DebugMode)
	}

	if server.AppEnv == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	// programmatically set swagger info
	router := gin.New()
	router.Use(ThisRouter.GinFunc)
	router.Use(gin.Recovery())
	router.Use(cors.New(configs.Cors))
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router = Routing(router)

	ThisRouter.Router = router
}

// InitTracing ..
func (ThisRouter *ThisRouter) InitTracing() {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "PROD"
	}

	tracer, reporter, closer, err := tracing.InitTracing(fmt.Sprintf("%s::%s", server.ServiceName, hostName), server.OpenTracingHost, tracing.WithEnableInfoLog(true))
	if err != nil {
		fmt.Println("Error :", err)
	}
	opentracing.SetGlobalTracer(tracer)

	ThisRouter.Closer = closer
	ThisRouter.Tracer = tracer
	ThisRouter.Err = err
	ThisRouter.Reporter = reporter
	ThisRouter.GinFunc = tracing.OpenTracer([]byte("api-request-"))
}

// Close ..
func (ThisRouter *ThisRouter) Close() {
	ThisRouter.Closer.Close()
}
