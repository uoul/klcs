package services

import (
	"maps"
	"slices"
	"sync"

	"github.com/uoul/go-async"
	"github.com/uoul/klcs/backend/core/domain"

	appError "github.com/uoul/klcs/backend/core/error"
)

// -------------------------------------------------------------------------------
// Type
// -------------------------------------------------------------------------------
type PrintService struct {
	mux      sync.RWMutex
	printers map[string][]async.Sequence[domain.PrintJob]

	printQueueSize uint
}

// -------------------------------------------------------------------------------
// Public
// -------------------------------------------------------------------------------

func (ps *PrintService) PrintJob(printerId string, job domain.PrintJob) error {
	// Lock for job sending
	ps.mux.RLock()
	defer ps.mux.RUnlock()
	// Check if printer connected
	client, exists := ps.printers[printerId]
	if !exists {
		return appError.NewErrNotFound("printer(%s) not connected", printerId)
	}
	// Send job to all printers using printerId
	for _, printer := range client {
		printer <- async.Success(job)
	}
	return nil
}

func (ps *PrintService) GetConnectedPrinters() []string {
	return slices.Collect(maps.Keys(ps.printers))
}

func (ps *PrintService) Register(printerId string) async.Sequence[domain.PrintJob] {
	// Lock for adding new printer
	ps.mux.Lock()
	defer ps.mux.Unlock()
	// Create new subscription
	sub := make(async.Sequence[domain.PrintJob], ps.printQueueSize)
	// Check if printer already registered
	_, exists := ps.printers[printerId]
	if exists {
		ps.printers[printerId] = append(ps.printers[printerId], sub)
	} else {
		ps.printers[printerId] = []async.Sequence[domain.PrintJob]{sub}
	}
	// Return new subscription
	return sub
}

func (ps *PrintService) UnRegister(printerId string, sub async.Sequence[domain.PrintJob]) {
	ps.mux.Lock()
	defer ps.mux.Unlock()
	// Check if printer is registered
	subs, exists := ps.printers[printerId]
	if !exists {
		return
	}
	// Find and remove the subscription
	for i, s := range subs {
		if s == sub {
			close(sub)
			subs[i] = nil // Set to nil --> Avoid memory leak
			subs = append(subs[:i], subs[i+1:]...)
			break
		}
	}
	// Update or remove printer entry
	if len(subs) == 0 {
		delete(ps.printers, printerId)
	} else {
		ps.printers[printerId] = subs
	}
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
		printers:       make(map[string][]async.Sequence[domain.PrintJob]),
		mux:            sync.RWMutex{},
		printQueueSize: 50,
	}
	for _, o := range opts {
		o(ps)
	}
	return ps
}
