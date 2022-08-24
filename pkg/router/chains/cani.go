package chains

import (
	"errors"
	"net/http"

	"github.com/fengxsong/vmdashboard/pkg/authorization/authorizer"
	"github.com/fengxsong/vmdashboard/pkg/endpoints/request"
	"github.com/fengxsong/vmdashboard/pkg/router/responsewriters"
)

func extractAttributesFromRequest(req *http.Request, getter authorizer.RequestAttributesGetter) (authorizer.Attributes, error) {
	ctx := req.Context()
	user, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.New("cannot find user info from context")
	}
	// use default getter here
	attr := getter.GetRequestAttributes(user, req)
	if attr == nil {
		return nil, errors.New("cannot find request info from context")
	}
	return attr, nil
}

// TODO: fix returning error message
func WithCanI(handler http.Handler, auth authorizer.Authorizer, getter authorizer.RequestAttributesGetter, failed http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		attr, err := extractAttributesFromRequest(req, getter)
		if err != nil {
			responsewriters.InternalError(w, req, errors.New("cannot get request attrs"))
			return
		}
		// TODO: return reason to failed handler
		decision, _, err := auth.Authorize(req.Context(), attr)
		if err != nil {
			responsewriters.InternalError(w, req, err)
			return
		}
		if decision == authorizer.DecisionDeny {
			// no need to return custom error
			failed.ServeHTTP(w, req)
			return
		}
		handler.ServeHTTP(w, req)
	})
}
