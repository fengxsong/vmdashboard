package rbac

import (
	"context"

	v1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	authzclient "k8s.io/client-go/kubernetes/typed/authorization/v1"

	"github.com/fengxsong/vmdashboard/pkg/authorization/authorizer"
)

type rbac struct {
	authz authzclient.AuthorizationV1Interface
}

func New(authz authzclient.AuthorizationV1Interface) *rbac {
	return &rbac{authz: authz}
}

func (r *rbac) Authorize(ctx context.Context, a authorizer.Attributes) (authorizer.Decision, string, error) {
	// no need for impersonate. lol
	ssar, err := r.authz.SubjectAccessReviews().Create(ctx, &v1.SubjectAccessReview{
		Spec: v1.SubjectAccessReviewSpec{
			User:   a.GetUser().GetName(),
			Groups: a.GetUser().GetGroups(),
			// todo: convert string slice into v1.ExtraValue type
			// Extra:  a.GetUser().GetExtra(),
			UID: a.GetUser().GetUID(),
			ResourceAttributes: &v1.ResourceAttributes{
				Namespace:   a.GetNamespace(),
				Verb:        a.GetVerb(),
				Group:       a.GetAPIGroup(),
				Version:     a.GetAPIVersion(),
				Resource:    a.GetResource(),
				Subresource: a.GetSubresource(),
				Name:        a.GetName(),
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return authorizer.DecisionDeny, "cannot create subjectaccessreview", err
	}
	if !ssar.Status.Allowed {
		return authorizer.DecisionDeny, "not allowed", nil
	}
	return authorizer.DecisionAllow, "", nil
}
