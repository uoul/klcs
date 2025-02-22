package services

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/hennedo/escpos"
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
	connectPrinter  func() (io.Writer, func() error, error)

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
	// Create printer
	p := escpos.New(conn)
	_, err = p.Initialize()
	if err != nil {
		return appError.NewErrPrint("failed to initialize printer - %v", err)
	}
	// Write Job name
	p.Size(2, 2).Justify(escpos.JustifyCenter).Underline(2)
	_, err = p.Write(job.ShopName)
	if err != nil {
		return appError.NewErrPrint("%v", err)
	}
	_, err = p.Underline(0).Size(2, 1).Justify(escpos.JustifyLeft).LineFeedD(2)
	if err != nil {
		return appError.NewErrPrint("%v", err)
	}
	// Write Articles
	for article, count := range job.OrderPositions {
		_, err = p.Write(fmt.Sprintf("%dx - %s", count, article))
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
		_, err = p.LineFeedD(2)
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
	}
	// Write Account if exists
	if len(job.AccountHolderName) > 0 {
		_, err = p.Write("------------------------")
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
		_, err = p.LineFeed()
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
		p.Justify(escpos.JustifyCenter)
		_, err = p.Write(fmt.Sprintf("Account: %s", job.AccountHolderName))
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
		_, err = p.LineFeed()
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
	}
	if len(job.Description) > 0 {
		_, err = p.Write("------------------------")
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
		_, err = p.LineFeed()
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
		p.Size(1, 1).Justify(escpos.JustifyCenter)
		_, err = p.Write(job.Description)
		if err != nil {
			return appError.NewErrPrint("%v", err)
		}
	}
	// Print
	err = p.PrintAndCut()
	if err != nil {
		return appError.NewErrPrint("%v", err)
	}
	return nil
}

// -------------------------------------------------------------------------------
// Options
// -------------------------------------------------------------------------------

func WithTcpConnector(addr string) func(*PrintService) {
	return func(ps *PrintService) {
		ps.connectPrinter = func() (io.Writer, func() error, error) {
			conn, err := net.Dial("tcp", addr)
			return conn, conn.Close, err
		}
	}
}

func WithUsbConnector(usbDev string) func(*PrintService) {
	return func(ps *PrintService) {
		ps.connectPrinter = func() (io.Writer, func() error, error) {
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
