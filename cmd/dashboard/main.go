package main

import (
	"fmt"
	"os"

	commonversion "github.com/prometheus/common/version"
	"github.com/spf13/pflag"

	"github.com/fengxsong/vmdashboard/cmd/dashboard/options"
	"github.com/fengxsong/vmdashboard/pkg/signals"
)

const app = "vmdashboard"

func main() {
	fs := pflag.NewFlagSet(app, pflag.ExitOnError)
	o := options.New()
	o.AddFlags(fs)

	if o.PrintVersion {
		fmt.Println(commonversion.Print(app))
		return
	}
	ctx := signals.SetupSignalHandler()
	server, err := o.Complete()
	if err != nil {
		fatal(err)
	}
	if err = server.Run(ctx); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}
