package main

import (
	"flag"
	"time"

	"github.com/uoul/go-common/log"
	"github.com/uoul/go-common/resource"
	"github.com/uoul/klcs/backend/print-agent/dal"
	"github.com/uoul/klcs/backend/print-agent/services"
)

func main() {
	// Arguments
	logLevel := flag.String("logLevel", "INFO", "OFF, TRACE, DEBUG, INFO, WARNING, ERROR, FATAL")
	klcsHost := flag.String("klcsHost", "", "BackedHost, where printjobs can be obtained (e.g. https://klcs.xxxx.xx)")
	klcsPrinterId := flag.String("klcsPrinterId", "", "PrinterId of klcs")
	timeZone := flag.String("timeZone", "Europe/Vienna", "Timezone")
	printerUsbAddr := flag.String("printerUsbAddr", "", "Path to usb device of printer (e.g. /dev/usb/lp3)")
	printerNetAddr := flag.String("printerNetAddr", "", "ip + port of printer (e.g. 192.168.0.10:9100)")
	flag.Parse()

	// Create Logger
	logger := log.NewConsoleLogger(log.StringToLogLevel(*logLevel, log.INFO))

	// Create ResourceManager
	rm := resource.NewResourceManager(10*time.Second, logger)

	// Create KlcsApi
	klcsApi := dal.NewKlcsApi(*klcsHost)

	// Create KlcsClientService
	klcsClient := services.NewKlcsClientService(logger, *klcsHost, *klcsPrinterId)
	rm.Register(klcsClient)
	defer rm.Unregister(klcsClient)

	// Create PrintService
	ts, err := time.LoadLocation(*timeZone)
	if err != nil {
		logger.Fatalf("Given timezone(%s) does not match IANA TimeZone format - %v", *timeZone, err)
		return
	}
	var printService services.IService
	if len(*printerNetAddr) > 0 {
		printService = services.NewPrintService(
			logger,
			klcsApi,
			klcsClient,
			*klcsPrinterId,
			ts,
			services.WithTcpConnector(*printerNetAddr),
		)
	} else if len(*printerUsbAddr) > 0 {
		printService = services.NewPrintService(
			logger,
			klcsApi,
			klcsClient,
			*klcsPrinterId,
			ts,
			services.WithUsbConnector(*printerUsbAddr),
		)
	} else {
		logger.Fatalf("Printer interface must be specified")
		return
	}
	rm.Register(printService)
	defer rm.Unregister(printService)

	// Run services
	go klcsClient.Run()
	go printService.Run()

	// Wait until termination
	rm.Wait()
}
