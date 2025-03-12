package services

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/uoul/escpos/netum/ns8360l"
	"github.com/uoul/go-common/log"
	"github.com/uoul/klcs/backend/print-agent/domain"
	appError "github.com/uoul/klcs/backend/print-agent/error"
)

// -------------------------------------------------------------------------------
// Type
// -------------------------------------------------------------------------------
type PrintService struct {
	logger      log.ILogger
	printJobSrc INotificationService[domain.PrintJob]

	printBufferSize int
	connectPrinter  func() (io.ReadWriter, func() error, error)

	stop chan any
}

// -------------------------------------------------------------------------------
// Public
// -------------------------------------------------------------------------------

// Close implements IService.
func (p *PrintService) Close() error {
	p.stop <- true
	return nil
}

// Run implements IService.
func (p *PrintService) Run() {
	// Subscribe to klcs api
	jobs := p.printJobSrc.Subscribe(uint(p.printBufferSize))
	defer p.printJobSrc.Unsubscribe(jobs)
	// Run
LP1:
	for {
		select {
		case job := <-jobs:
			if job.Error == nil {
				err := p.printOrder(&job.Result)
				if err != nil {
					p.logger.Errorf("%v", err)
				}
			}
		case <-p.stop:
			break LP1
		}
	}
}

// -------------------------------------------------------------------------------
// Private
// -------------------------------------------------------------------------------
func (ps *PrintService) printOrder(job *domain.PrintJob) error {
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
	if err := p.Print(
		time.Now().Format("02.01.2006 15:04:05\n"),
		p.WithJustifyCenter(),
	); err != nil {
		return appError.NewErrPrint("failed to print timestamp - %v", err)
	}
	// Cut
	return p.Cut()
}

// -------------------------------------------------------------------------------
// Options
// -------------------------------------------------------------------------------

func WithTcpConnector(addr string) func(*PrintService) {
	return func(ps *PrintService) {
		ps.connectPrinter = func() (io.ReadWriter, func() error, error) {
			conn, err := net.Dial("tcp", addr)
			return conn, conn.Close, err
		}
	}
}

func WithUsbConnector(usbDev string) func(*PrintService) {
	return func(ps *PrintService) {
		ps.connectPrinter = func() (io.ReadWriter, func() error, error) {
			f, err := os.OpenFile(usbDev, os.O_RDWR, 0)
			return f, f.Close, err
		}
	}
}

func WithPrintBufferSize(size int) func(*PrintService) {
	return func(ps *PrintService) {
		ps.printBufferSize = size
	}
}

// ----------------------------------------------------------------------
// Constructor
// ----------------------------------------------------------------------

func NewPrintService(logger log.ILogger, klcsClient INotificationService[domain.PrintJob], opts ...func(*PrintService)) IService {
	ps := &PrintService{
		logger:          logger,
		printJobSrc:     klcsClient,
		printBufferSize: 50,
		connectPrinter:  nil,
		stop:            make(chan any, 1),
	}
	for _, o := range opts {
		o(ps)
	}
	return ps
}
