package cmd

import (
	"fmt"
	"os"

	flog "github.com/everywan/foundation-go/log"

	"github.com/everywan/identifier"
	"github.com/everywan/identifier/internal/services"
)

type bootstrap struct {
	Opts     *ApplicationOps
	Logger   *flog.Logger
	SfSvc    identifier.SnowflakeService
	Teardown func()
}

func newBootstrap(opts *ApplicationOps) (boot *bootstrap, err error) {
	boot = new(bootstrap)
	logger := flog.NewLogger(opts.Logger, os.Stdout)
	sfSvc, err := services.NewSnowflakeService(opts.Snowflake.WorkerID)
	if err != nil {
		return boot, err
	}
	teardown := func() { fmt.Println("teardown bootstrap!") }
	boot = &bootstrap{
		Opts:     opts,
		Logger:   logger,
		SfSvc:    sfSvc,
		Teardown: teardown,
	}
	return boot, nil
}
