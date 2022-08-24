package responsewriters

import (
	"fmt"
	"net/http"
	"strings"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

type Func func(http.ResponseWriter, *http.Request, error)

// Avoid emitting errors that look like valid HTML. Quotes are okay.
var sanitizer = strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;")

// InternalError renders a simple internal error
func InternalError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, sanitizer.Replace(fmt.Sprintf("Internal Server Error: %q: %v", req.RequestURI, err)),
		http.StatusInternalServerError)
	utilruntime.HandleError(err)
}

func Unauthorized(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, sanitizer.Replace(fmt.Sprintf("%q: %v", req.RequestURI, err)), http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, sanitizer.Replace(fmt.Sprintf("%q: %v", req.RequestURI, err)), http.StatusForbidden)
}
