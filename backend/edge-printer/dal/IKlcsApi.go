package dal

import (
	"github.com/uoul/go-common/async"
	"github.com/uoul/klcs/backend/edge-printer/domain"
)

type IKlcsApi interface {
	GetPrintJobs() async.Stream[domain.PrintJob]
}
