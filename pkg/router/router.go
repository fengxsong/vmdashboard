package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"

	"github.com/fengxsong/vmdashboard/pkg/csrf"
)

type Server struct {
	option     *Option
	restConfig *rest.Config
	logger     *zap.SugaredLogger

	handler http.Handler
}

func New(o *Option, cfg *rest.Config, logger *zap.SugaredLogger) *Server {
	return &Server{
		option:     o,
		restConfig: cfg,
		logger:     logger,
	}
}

func (s *Server) registerAPIs() error {
	s.logger.Info("Adding Kube APIs")
	apiProxyPrefix := DefaultLegacyAPIPrefix + "/"
	apisProxyPrefix := APIGroupPrefix + "/"
	proxyHandlerAPI, err := NewProxyHandler(apiProxyPrefix, s.restConfig, s.option.proxyKeepAliveTimeout)
	if err != nil {
		return err
	}
	proxyHandlerAPIs, err := NewProxyHandler(apisProxyPrefix, s.restConfig, s.option.proxyKeepAliveTimeout)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	o, err := s.option.Complete(s.restConfig)
	if err != nil {
		return err
	}
	wrapProxy := func(handler http.Handler, identity string) http.Handler {
		return defaultBuildHandlerChain(handler, o, s.logger.With("router", identity))
	}
	mux.Handle(apiProxyPrefix, wrapProxy(proxyHandlerAPI, DefaultLegacyAPIPrefix))
	mux.Handle(apisProxyPrefix, wrapProxy(proxyHandlerAPIs, APIGroupPrefix))

	// static resource endpoint
	mux.Handle(s.option.staticPrefix, http.StripPrefix(s.option.staticPrefix, staticEndpoint(s.option.staticRootDir)))

	// TODO: add some health checker
	health := healthcheck.NewHandler()
	mux.HandleFunc("/healthz/live", health.LiveEndpoint)
	mux.HandleFunc("/healthz/ready", health.ReadyEndpoint)

	// TODO: extend some other APIs
	// /auth/token

	s.handler = mux
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.registerAPIs(); err != nil {
		return err
	}

	CSRF := csrf.Protect()
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.option.port),
		Handler: CSRF(s.handler),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.logger.Info("Shuting down")
		return server.Shutdown(context.Background())
	}
}
