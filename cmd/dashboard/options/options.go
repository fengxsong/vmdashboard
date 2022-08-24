package options

import (
	"context"
	"os"

	"github.com/spf13/pflag"

	"github.com/fengxsong/vmdashboard/pkg/environment"
	"github.com/fengxsong/vmdashboard/pkg/log/zap"
	"github.com/fengxsong/vmdashboard/pkg/router"
)

type Runnable interface {
	Run(context.Context) error
}

type ServerRunOption struct {
	zo           *zap.Options
	env          *environment.ClientConfig
	ro           *router.Option
	PrintVersion bool
}

func New() *ServerRunOption {
	return &ServerRunOption{
		zo:  &zap.Options{},
		env: &environment.ClientConfig{},
		ro:  &router.Option{},
	}
}

func (o *ServerRunOption) AddFlags(fs *pflag.FlagSet) {
	o.zo.BindFlags(fs)
	o.env.InitFlags(fs)
	o.ro.AddFlags(fs)
	fs.BoolVarP(&o.PrintVersion, "version", "v", false, "Print version information")

	fs.Parse(os.Args)
}

func (o *ServerRunOption) Complete() (Runnable, error) {
	cfg, err := o.env.GetRESTConfig()
	if err != nil {
		return nil, err
	}
	logger := zap.New(zap.UseFlagOptions(o.zo))
	return router.New(o.ro, cfg, logger), nil
}
