package router

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"

	"github.com/fengxsong/vmdashboard/pkg/endpoints/request"
)

const (
	// DefaultLegacyAPIPrefix is where the legacy APIs will be located.
	DefaultLegacyAPIPrefix = "/api"

	// APIGroupPrefix is where non-legacy API group will be located.
	APIGroupPrefix = "/apis"
)

type responder struct {
	logger *zap.SugaredLogger
}

func (r *responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	r.logger.Errorf("Error while proxying request: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func makeUpgradeTransport(config *rest.Config, keepalive time.Duration) (proxy.UpgradeRequestRoundTripper, error) {
	transportConfig, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	tlsConfig, err := transport.TLSConfigFor(transportConfig)
	if err != nil {
		return nil, err
	}
	rt := utilnet.SetOldTransportDefaults(&http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: keepalive,
		}).DialContext,
	})

	upgrader, err := transport.HTTPWrappersForConfig(transportConfig, proxy.MirrorRequest)
	if err != nil {
		return nil, err
	}
	return proxy.NewUpgradeRequestRoundTripper(rt, upgrader), nil
}

// NewProxyHandler creates an API proxy handler for the cluster
func NewProxyHandler(apiProxyPrefix string, cfg *rest.Config, keepalive time.Duration) (http.Handler, error) {
	host := cfg.Host
	if !strings.HasSuffix(host, "/") {
		host = host + "/"
	}
	target, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	responder := &responder{}
	transport, err := rest.TransportFor(cfg)
	if err != nil {
		return nil, err
	}
	transport = &wrappedRoundTripper{transport, true}
	upgradeTransport, err := makeUpgradeTransport(cfg, keepalive)
	if err != nil {
		return nil, err
	}
	proxy := proxy.NewUpgradeAwareHandler(target, transport, false, false, responder)
	proxy.UpgradeTransport = upgradeTransport
	proxy.UseRequestLocation = true
	proxy.UseLocationHost = true

	proxyServer := http.Handler(proxy)

	return proxyServer, nil
}

type wrappedRoundTripper struct {
	inner       http.RoundTripper
	impersonate bool
}

func (rt *wrappedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !rt.impersonate {
		return rt.inner.RoundTrip(req)
	}
	user, ok := request.UserFrom(req.Context())
	if !ok || user == nil || len(req.Header.Get(transport.ImpersonateUserHeader)) != 0 {
		return rt.inner.RoundTrip(req)
	}
	req = utilnet.CloneRequest(req)
	if username := user.GetName(); username != "" {
		req.Header.Set(transport.ImpersonateUserHeader, user.GetName())
	}
	if uid := user.GetUID(); uid != "" {
		req.Header.Set(transport.ImpersonateUIDHeader, uid)
	}
	for _, group := range user.GetGroups() {
		req.Header.Add(transport.ImpersonateGroupHeader, group)
	}
	for k, vv := range user.GetExtra() {
		for _, v := range vv {
			req.Header.Add(transport.ImpersonateUserExtraHeaderPrefix+http.CanonicalHeaderKey(k), v)
		}
	}
	return rt.inner.RoundTrip(req)
}
