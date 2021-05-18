package internal

import (
	"github.com/sha1n/benchy/api"
	log "github.com/sirupsen/logrus"
)

// WriteReportFnFor returns a wrapper function for the specified api.WriteReportFn
func WriteReportFnFor(write api.WriteReportFn) api.WriteReportFn {
	w := func(ts api.Summary, config *api.BenchmarkSpec, ctx *api.ReportContext) (err error) {
		log.Info("Writing report...")
		defer log.Info("Done!")

		return write(ts, config, ctx)
	}

	return w
}
