package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/fengxsong/vmdashboard/pkg/authentication/authenticator"
	"github.com/fengxsong/vmdashboard/pkg/authentication/user"
)

var (
	errInvalidToken = errors.New("token is invalid")
	errExpiredToken = errors.New("token has expired")
)

type token struct {
	secret []byte
	period time.Duration
}

type Option func(*token)

func WithSecret(secret []byte) Option {
	return func(t *token) {
		t.secret = secret
	}
}

func WithPeriod(d time.Duration) Option {
	return func(t *token) {
		t.period = d
	}
}

func (t *token) AuthenticateToken(ctx context.Context, token string) (*authenticator.Response, bool, error) {
	cl, err := t.verifyToken(token)
	if err != nil {
		return nil, false, err
	}
	return &authenticator.Response{
		User: &user.DefaultInfo{
			Name:   cl.Name,
			UID:    cl.UID,
			Groups: cl.Groups,
			Extra:  cl.Extra,
		},
	}, true, nil
}

func (t *token) Issue(u user.Info) (string, error) {
	cl := &claim{
		Name:   u.GetName(),
		UID:    u.GetUID(),
		Groups: u.GetGroups(),
		Extra:  u.GetExtra(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(t.period).Unix(),
		},
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	tkStr, err := tk.SignedString(t.secret)
	return tkStr, err
}

type claim struct {
	Name   string              `json:"name"`
	UID    string              `json:"uid,omitempty"`
	Groups []string            `json:"groups,omitempty"`
	Extra  map[string][]string `json:"extra,omitempty"`
	jwt.StandardClaims
}

func (t *token) verifyToken(token string) (*claim, error) {
	keyFunc := func(tk *jwt.Token) (interface{}, error) {
		_, ok := tk.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errInvalidToken
		}
		return t.secret, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &claim{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errExpiredToken) {
			return nil, errExpiredToken
		}
		return nil, errInvalidToken
	}

	cl, ok := jwtToken.Claims.(*claim)
	if !ok {
		return nil, errInvalidToken
	}

	return cl, nil
}

func New(opts ...Option) *token {
	t := &token{}
	for _, o := range opts {
		o(t)
	}
	return t
}
