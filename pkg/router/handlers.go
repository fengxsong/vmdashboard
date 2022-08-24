package router

import (
	"net/http"

	"github.com/fengxsong/vmdashboard/dist"
	"github.com/fengxsong/vmdashboard/pkg/authentication/user"
)

func staticEndpoint(rootDir string) http.Handler {
	var staticFs http.FileSystem
	if len(rootDir) == 0 {
		staticFs = http.FS(dist.StaticFS)
	} else {
		staticFs = http.Dir(rootDir)
	}
	fs := http.FileServer(staticFs)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// todo: add cache-control options
		fs.ServeHTTP(rw, req)
	})
}

type tokenIssuer interface {
	Issue(user.Info) (string, error)
}

// TODO: auth handler
// builtin/ldap/oidc
