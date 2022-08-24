package router

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fengxsong/vmdashboard/pkg/authentication/authenticator"
	anonymousauthentication "github.com/fengxsong/vmdashboard/pkg/authentication/request/anonymous"
	"github.com/fengxsong/vmdashboard/pkg/authentication/request/bearertoken"
	"github.com/fengxsong/vmdashboard/pkg/authentication/token/jwt"
	"github.com/fengxsong/vmdashboard/pkg/authentication/token/tokenreview"
	anonymousauthorization "github.com/fengxsong/vmdashboard/pkg/authorization/anonymous"
	"github.com/fengxsong/vmdashboard/pkg/authorization/authorizer"
	"github.com/fengxsong/vmdashboard/pkg/authorization/rbac"
)

type Option struct {
	port                  int
	proxyKeepAliveTimeout time.Duration
	tokenPeriod           time.Duration
	jwtSecretKey          string
	staticPrefix          string
	staticRootDir         string
	authType              string
}

func (o *Option) AddFlags(fs *pflag.FlagSet) {
	fs.IntVarP(&o.port, "port", "p", 8080, "HTTP Listening port")
	fs.DurationVar(&o.proxyKeepAliveTimeout, "keepalive", 30*time.Second, "Proxy keepalive timeout")
	fs.DurationVar(&o.tokenPeriod, "token-period", 2*time.Hour, "JWT token valid period")
	fs.StringVar(&o.jwtSecretKey, "jwt-secret", "", "JWT secret key")
	fs.StringVar(&o.staticPrefix, "static-prefix", "/static", "URL prefix for static resource")
	fs.StringVar(&o.staticRootDir, "static-rootdir", "", "Rootdir for static resource, if not specified then it will use embed filesystem for **dist** dir")
	fs.StringVar(&o.authType, "auth-type", "off", "Token auth type, avaliable option are off, jwt, tokenreview")
}

// TODO: make this function more flexible
func (o *Option) Complete(cfg *rest.Config) (*completeOption, error) {
	if !strings.HasSuffix(o.staticPrefix, "/") {
		o.staticPrefix += "/"
	}
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	ctrlclient, err := client.New(cfg, client.Options{})
	if err != nil {
		return nil, err
	}
	co := &completeOption{
		cs: cs,
	}

	if o.authType == "off" {
		co.authenticator = anonymousauthentication.NewAuthenticator()
		co.authorizer = anonymousauthorization.New()
	} else {
		type token interface {
			authenticator.Token
			tokenIssuer
		}
		var tk token
		switch o.authType {
		case "jwt":
			tk = jwt.New(jwt.WithPeriod(o.tokenPeriod), jwt.WithSecret([]byte(o.jwtSecretKey)))
		case "tokenreview":
			tk = tokenreview.New(co.cs.AuthenticationV1(), ctrlclient)
		default:
			return nil, fmt.Errorf("unsupported auth type %s", o.authType)
		}
		co.authenticator = bearertoken.New(tk)
		co.authorizer = rbac.New(co.cs.AuthorizationV1())
		co.tokenIssuer = tk
	}

	return co, nil
}

type completeOption struct {
	// *Option
	cs            kubernetes.Interface
	authenticator authenticator.Request
	authorizer    authorizer.Authorizer
	tokenIssuer   tokenIssuer
}
