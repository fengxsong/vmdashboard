package chains

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/fengxsong/vmdashboard/pkg/endpoints/request"
	"github.com/fengxsong/vmdashboard/pkg/router/responsewriters"
)

// WithRequestInfo attaches a RequestInfo to the context.
func WithRequestInfo(handler http.Handler, resolver request.RequestInfoResolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		info, err := resolver.NewRequestInfo(req)
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("failed to create RequestInfo: %v", err))
			return
		}

		req = req.WithContext(request.WithRequestInfo(ctx, info))

		handler.ServeHTTP(w, req)
	})
}

// COPY from vendor/k8s.io/apiserver/pkg/server/config.go > NewRequestInfoResolver(*Config)
func NewRequestInfoResolver(apiGroupPrefix string, legacyAPIGroupPrefixes []string, logger *zap.SugaredLogger) *request.RequestInfoFactory {
	apiPrefixes := sets.NewString(strings.Trim(apiGroupPrefix, "/")) // all possible API prefixes
	legacyAPIPrefixes := sets.String{}                               // APIPrefixes that won't have groups (legacy)
	for _, legacyAPIPrefix := range legacyAPIGroupPrefixes {
		apiPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
		legacyAPIPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
	}

	return &request.RequestInfoFactory{
		APIPrefixes:          apiPrefixes,
		GrouplessAPIPrefixes: legacyAPIPrefixes,
	}
}
