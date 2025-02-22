package services

import (
	"sync"

	"github.com/uoul/go-common/async"
	"github.com/uoul/klcs/backend/core/domain"

	appError "github.com/uoul/klcs/backend/core/error"
)

// -------------------------------------------------------------------------------
// Type
// -------------------------------------------------------------------------------

type PrintService struct {
	mux     sync.Mutex
	clients map[string]async.Stream[domain.PrintJob]

	printQueueSize uint
}

// -------------------------------------------------------------------------------
// Public
// -------------------------------------------------------------------------------

func (ps *PrintService) PrintJob(printerId string, job domain.PrintJob) error {
	client, exists := ps.clients[printerId]
	if !exists {
		return appError.NewErrNotFound("printer(%s) not connected", printerId)
	}
	client <- async.ActionResult[domain.PrintJob]{
		Result: job,
		Error:  nil,
	}
	return nil
}

func (ps *PrintService) Subscribe(printerId string) (async.Stream[domain.PrintJob], error) {
	ps.mux.Lock()
	defer ps.mux.Unlock()
	_, exists := ps.clients[printerId]
	if exists {
		return nil, appError.NewErrConflict("printer already exists")
	}
	client := async.NewBufferedStream[domain.PrintJob](ps.printQueueSize)
	ps.clients[printerId] = client
	return client, nil
}

func (ps *PrintService) Unsubscribe(printerId string) {
	ps.mux.Lock()
	defer ps.mux.Unlock()
	delete(ps.clients, printerId)
}

// -------------------------------------------------------------------------------
// Private
// -------------------------------------------------------------------------------

// -------------------------------------------------------------------------------
// Options
// -------------------------------------------------------------------------------

func WithPrintBufferSize(size uint) func(*PrintService) {
	return func(ps *PrintService) {
		ps.printQueueSize = size
	}
}

// ----------------------------------------------------------------------
// Constructor
// ----------------------------------------------------------------------

func NewPrintService(opts ...func(*PrintService)) *PrintService {
	ps := &PrintService{
		clients:        map[string]async.Stream[domain.PrintJob]{},
		mux:            sync.Mutex{},
		printQueueSize: 50,
	}
	for _, o := range opts {
		o(ps)
	}
	return ps
}
