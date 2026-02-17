package dal

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/uoul/escpos/netum/ns8360l"
	"github.com/uoul/klcs/backend/print-agent/domain"
	appError "github.com/uoul/klcs/backend/print-agent/error"
)

// -------------------------------------------------------------------------------
// Type
// -------------------------------------------------------------------------------
type PrinterClient struct {
	timeLocation   *time.Location
	connectPrinter func() (io.ReadWriter, func() error, error)
}

// -------------------------------------------------------------------------------
// Public
// -------------------------------------------------------------------------------
func (ps *PrinterClient) PrintOrder(job domain.PrintJob) error {
	if ps.connectPrinter == nil {
		return appError.NewErrConfig("no printer connection type has been configured (WithTcpConnector, WithUsbConnector)")
	}
	conn, close, err := ps.connectPrinter()
	if err != nil {
		return appError.NewErrIO("failed to connect to printer - %v", err)
	}
	defer close()
	// Create Printer
	p := ns8360l.NewPrinter(conn)
	// Print Headline (ShopName)
	if err := p.Print(
		fmt.Sprintf("%s\n\n", job.ShopName),
		p.WithSize(3, 2),
		p.WithUnderline(2),
		p.WithJustifyCenter(),
	); err != nil {
		return appError.NewErrPrint("failed to print headline - %v", err)
	}
	// Print Articles
	for article, count := range job.OrderPositions {
		if err := p.Print(
			fmt.Sprintf("%dx - %s\n", count, article),
			p.WithSize(1, 2),
		); err != nil {
			return appError.NewErrPrint("failed to print order position - %v", err)
		}
	}
	// Print Description
	if len(job.Description) > 0 {
		if err := p.Print(
			"------------------------\n",
			p.WithSize(1, 2),
			p.WithJustifyCenter(),
		); err != nil {
			return appError.NewErrPrint("failed to print divider - %v", err)
		}
		if err := p.Print(
			fmt.Sprintf("%s\n", job.Description),
			p.WithJustifyCenter(),
		); err != nil {
			return appError.NewErrPrint("failed to print description")
		}
	}
	// Print Account if exists
	if len(job.AccountHolderName) > 0 {
		if err := p.Print(
			"------------------------\n",
			p.WithSize(1, 2),
			p.WithJustifyCenter(),
		); err != nil {
			return appError.NewErrPrint("failed to print divider - %v", err)
		}
		if err := p.Print(
			fmt.Sprintf("Account: %s\n", job.AccountHolderName),
			p.WithJustifyCenter(),
		); err != nil {
			return appError.NewErrPrint("failed to print account holder - %v", err)
		}
	}
	// Print Timestamp
	if err := p.Print(
		"------------------------\n",
		p.WithSize(1, 2),
		p.WithJustifyCenter(),
	); err != nil {
		return appError.NewErrPrint("failed to print divider - %v", err)
	}
	dt, err := time.Parse(time.RFC3339, job.Timestamp)
	if err != nil {
		return appError.NewErrPrint("failed to parse timestamp - %v", err)
	}
	if err := p.Print(
		dt.In(ps.timeLocation).Format("02.01.2006 15:04:05\n"),
		p.WithJustifyCenter(),
	); err != nil {
		return appError.NewErrPrint("failed to print timestamp - %v", err)
	}
	// Cut
	if err := p.Cut(); err != nil {
		return appError.NewErrPrint("failed to cut paper - %v", err)
	}
	// Success
	return nil
}

// -------------------------------------------------------------------------------
// Options
// -------------------------------------------------------------------------------

func WithTcpConnector(addr string) func(*PrinterClient) {
	return func(ps *PrinterClient) {
		ps.connectPrinter = func() (io.ReadWriter, func() error, error) {
			conn, err := net.Dial("tcp", addr)
			return conn, conn.Close, err
		}
	}
}

func WithUsbConnector(usbDev string) func(*PrinterClient) {
	return func(ps *PrinterClient) {
		ps.connectPrinter = func() (io.ReadWriter, func() error, error) {
			f, err := os.OpenFile(usbDev, os.O_RDWR, 0)
			return f, f.Close, err
		}
	}
}

// ----------------------------------------------------------------------
// Constructor
// ----------------------------------------------------------------------

func NewPrinterClient(timeLocation *time.Location, opts ...func(*PrinterClient)) *PrinterClient {
	ps := &PrinterClient{
		timeLocation:   timeLocation,
		connectPrinter: nil,
	}
	for _, o := range opts {
		o(ps)
	}
	return ps
}
