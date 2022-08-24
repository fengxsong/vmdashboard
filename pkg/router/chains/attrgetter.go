package chains

import (
	"net/http"

	"github.com/fengxsong/vmdashboard/pkg/authentication/user"
	"github.com/fengxsong/vmdashboard/pkg/authorization/authorizer"
	"github.com/fengxsong/vmdashboard/pkg/endpoints/request"
)

type GetterFunc func(user.Info, *http.Request) authorizer.Attributes

func (f GetterFunc) GetRequestAttributes(u user.Info, req *http.Request) authorizer.Attributes {
	return f(u, req)
}

func DefaultRequestAttributesGetter(u user.Info, req *http.Request) authorizer.Attributes {
	ri, ok := request.RequestInfoFrom(req.Context())
	if !ok {
		return nil
	}
	return authorizer.AttributesRecord{
		User:            u,
		Verb:            ri.Verb,
		Namespace:       ri.Namespace,
		APIGroup:        ri.APIGroup,
		APIVersion:      ri.APIVersion,
		Resource:        ri.Resource,
		Subresource:     ri.Subresource,
		Name:            ri.Name,
		ResourceRequest: ri.IsResourceRequest,
		Path:            ri.Path,
	}
}
