package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	gonfig "github.com/uoul/go-config"
	"github.com/uoul/klcs/backend/print-agent/config"
	"github.com/uoul/klcs/backend/print-agent/dal"
	"github.com/uoul/klcs/backend/print-agent/services"
)

const (
	VERSION = "{VERSION}"
)

func main() {
	// Create Application Context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Load Config
	cfg, err := gonfig.FromFlags[config.AppConfig]()
	if err != nil {
		panic(err)
	}
	// Setup Logger
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLvl(),
	})
	slog.SetDefault(slog.New(logHandler))
	// Log current version
	slog.Info("Starting print-agent", slog.String("version", VERSION))
	// Create KlcsApi
	klcsApi := dal.NewKlcsApi(cfg.KlcsHost)
	// Create PrinterClient
	tz, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		slog.Error("Given timezone does not match IANA TimeZone format", slog.String("timezone", cfg.TimeZone), slog.Any("error", err))
		os.Exit(1)
	}
	var printerClient *dal.PrinterClient
	if len(cfg.Printer.NetAddr) > 0 {
		printerClient = dal.NewPrinterClient(
			tz,
			dal.WithTcpConnector(cfg.Printer.NetAddr),
		)
	} else if len(cfg.Printer.UsbAddr) > 0 {
		printerClient = dal.NewPrinterClient(
			tz,
			dal.WithUsbConnector(cfg.Printer.UsbAddr),
		)
	} else {
		slog.Error("No printer interface specified")
		os.Exit(1)
	}
	// Run Services
	services.NewPrintService(ctx, cfg.Printer.Id, klcsApi, printerClient)
	// Wait until termination

	<-ctx.Done()
}
