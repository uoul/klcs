package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/uoul/klcs/backend/print-agent/dal"
)

// -------------------------------------------------------------------------------
// Type
// -------------------------------------------------------------------------------

type PrintService struct {
	printerClient *dal.PrinterClient
	klcsClient    *dal.KlcsApi
	printerId     string

	retryCooldown time.Duration
}

// -------------------------------------------------------------------------------
// Public
// -------------------------------------------------------------------------------

// -------------------------------------------------------------------------------
// Private
// -------------------------------------------------------------------------------

func mustJson(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (k *PrintService) run(ctx context.Context) error {
	// Get PrintJobs
	jobs, err := k.klcsClient.GetPrintJobs(ctx, k.printerId)
	if err != nil {
		return err
	}
	// Listen for jobs
	for {
		select {
		case <-ctx.Done():
			return nil // Shutdown
		case job := <-jobs:
			// Check if error
			if job.Error != nil {
				return job.Error
			}
			slog.Debug("new incomming printjob", slog.String("job", mustJson(job.Value)))
			// Print Job
			if err := k.printerClient.PrintOrder(job.Value); err != nil {
				return err
			}
			// Acknowledge
			if err := k.klcsClient.AcknowledgePrintJob(ctx, k.printerId, job.Value.TransactionId); err != nil {
				return err
			}
			slog.Info("Job completed successfully", slog.String("job", mustJson(job.Value)))
		}
	}
}

// -------------------------------------------------------------------------------
// Options
// -------------------------------------------------------------------------------

func WithPrintServiceRetryCooldown(cooldown time.Duration) func(*PrintService) {
	return func(kps *PrintService) {
		kps.retryCooldown = cooldown
	}
}

// ----------------------------------------------------------------------
// Constructor
// ----------------------------------------------------------------------

func NewPrintService(ctx context.Context, printerId string, klcsApi *dal.KlcsApi, printerClient *dal.PrinterClient, opts ...func(*PrintService)) *PrintService {
	s := &PrintService{
		klcsClient:    klcsApi,
		printerClient: printerClient,
		printerId:     printerId,

		retryCooldown: 30 * time.Second,
	}
	for _, o := range opts {
		o(s)
	}
	// Run Service
	go func() {
		for {
			err := s.run(ctx)
			if err == nil {
				break // Shutdown
			}
			slog.Error("KlcsPrinter service failed", slog.Any("error", err))
			time.Sleep(s.retryCooldown)
		}
	}()
	return s
}
