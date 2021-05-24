package report

import "github.com/sha1n/benchy/api"

type RawDataHandler interface {
	Handle(api.Trace) error
}
