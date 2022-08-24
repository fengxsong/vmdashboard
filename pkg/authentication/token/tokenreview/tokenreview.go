package tokenreview

import (
	"context"
	"errors"

	authnv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	authnclient "k8s.io/client-go/kubernetes/typed/authentication/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fengxsong/vmdashboard/pkg/authentication/authenticator"
	"github.com/fengxsong/vmdashboard/pkg/authentication/user"
)

type auth struct {
	authn authnclient.AuthenticationV1Interface
	// controller-runtime cached client
	ctrlclient client.Client
}

func New(authn authnclient.AuthenticationV1Interface, ctrlclient client.Client) *auth {
	return &auth{
		authn:      authn,
		ctrlclient: ctrlclient,
	}
}

func (a auth) AuthenticateToken(ctx context.Context, token string) (*authenticator.Response, bool, error) {
	tr, err := a.authn.TokenReviews().Create(ctx, &authnv1.TokenReview{
		Spec: authnv1.TokenReviewSpec{
			Token: token,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}
	// should never happened :) when we use a token from another cluster
	if len(tr.Status.Error) > 0 {
		return nil, false, errors.New(tr.Status.Error)
	}

	if !tr.Status.Authenticated {
		return nil, false, nil
	}
	return &authenticator.Response{
		User: &user.DefaultInfo{
			Name:   tr.Status.User.Username,
			Groups: tr.Status.User.Groups,
			// Extra:  tr.Status.User.Extra,
			UID: tr.Status.User.UID,
		},
	}, true, nil
}

func (a auth) Issue(u user.Info) (string, error) {
	// TODO: maybe use cached client to look for token in sa's secret
	return "", nil
}
