package chains

import (
	"net/http"

	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/fengxsong/vmdashboard/pkg/authentication/authenticator"
	"github.com/fengxsong/vmdashboard/pkg/endpoints/request"
	"github.com/fengxsong/vmdashboard/pkg/router/responsewriters"
)

func WithAuthentication(handler http.Handler, auth authenticator.Request, failed http.Handler, logger *zap.SugaredLogger) http.Handler {
	if auth == nil {
		logger.Warn("Authentication is disabled")
		return handler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		resp, ok, err := auth.AuthenticateRequest(req)
		if err != nil || !ok {
			if err != nil {
				logger.Errorw("fail authenticate", "err", err)
			}
			req = req.WithContext(request.WithError(req.Context(), err))
			failed.ServeHTTP(w, req)
			return
		}
		req.Header.Del("Authorization")

		req = req.WithContext(request.WithUser(req.Context(), resp.User))

		handler.ServeHTTP(w, req)
	})
}

func FailHandler(rw responsewriters.Func, defaultError error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		err, found := request.ErrorFrom(ctx)
		if found {
			rw(w, req, err)
			return
		}
		rw(w, req, defaultError)
	})
}

func Unauthorized() http.Handler {
	return FailHandler(responsewriters.Unauthorized, apierrors.NewUnauthorized("Unauthorized"))
}

func Forbidden() http.Handler {
	return FailHandler(responsewriters.Forbidden, apierrors.NewForbidden(schema.GroupResource{}, "", nil))
}
