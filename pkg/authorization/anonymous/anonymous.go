package anonymous

import (
	"context"

	"github.com/fengxsong/vmdashboard/pkg/authorization/authorizer"
)

func New() authorizer.Authorizer {
	return authorizer.AuthorizerFunc(func(ctx context.Context, a authorizer.Attributes) (authorizer.Decision, string, error) {
		return authorizer.DecisionAllow, "", nil
	})
}
