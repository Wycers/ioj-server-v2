package http


import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/pkg/utils/netutil"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

const specialKey = "imf1nlTy0j"

type Options struct {
	Host string
	Port int
	Mode string
}

type Server struct {
	o          *Options
	app        string
	host       string
	port       int
	logger     *zap.Logger
	router     *gin.Engine
	httpServer http.Server
}

func NewOptions(v *viper.Viper) (*Options, error) {
	var (
		o   = new(Options)
		err error
	)

	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, err
	}

	return o, err
}

type InitControllers func(r *gin.Engine)

func NewRouter(o *Options, logger *zap.Logger, init InitControllers, tracer opentracing.Tracer) *gin.Engine {

	// 配置gin
	gin.SetMode(o.Mode)
	r := gin.New()

	r.Use(gin.Recovery()) // panic之后自动恢复
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	//r.Use(middlewares..New(r).Middleware()) // 添加 prometheus 监控
	r.Use(ginhttp.Middleware(tracer))
	store := cookie.NewStore([]byte(specialKey))
	r.Use(sessions.Sessions("ioj-session", store))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	pprof.Register(r)

	init(r)

	return r
}

func New(o *Options, logger *zap.Logger, router *gin.Engine) (*Server, error) {
	var s = &Server{
		logger: logger.With(zap.String("type", "http.Server")),
		router: router,
		o:      o,
	}

	return s, nil
}

func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	s.port = s.o.Port
	if s.port == 0 {
		s.port = netutil.GetAvailablePort()
	}

	s.host = s.o.Host
	if s.host == "" {
		s.host = "127.0.0.1"
	}

	if s.host == "" {
		return errors.New("get local ipv4 error")
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.httpServer = http.Server{Addr: addr, Handler: s.router}

	s.logger.Info("http server starting ...", zap.String("addr", addr))
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("start http server err", zap.Error(err))
			return
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}

	return nil
}

var ProviderSet = wire.NewSet(New, NewRouter, NewOptions)
