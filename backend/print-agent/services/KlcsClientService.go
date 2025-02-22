package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/log"
	"github.com/uoul/klcs/backend/print-agent/domain"
	appError "github.com/uoul/klcs/backend/print-agent/error"
)

const (
	READ_BUFFER_SIZE = 10240
)

// -------------------------------------------------------------------------------
// Type
// -------------------------------------------------------------------------------

type KlcsClientService struct {
	logger log.ILogger

	stop       chan any
	clients    map[async.Stream[domain.PrintJob]]bool
	httpClient http.Client

	klcsBackendHost string
	printerId       string
}

// -------------------------------------------------------------------------------
// Public
// -------------------------------------------------------------------------------

// Close implements INotificationService.
func (k *KlcsClientService) Close() error {
	k.stop <- true
	return nil
}

// Run implements INotificationService.
func (k *KlcsClientService) Run() {
LP1:
	for {
		select {
		case <-k.stop:
			break LP1
		default:
			err := k.connectKlcsEventStream()
			if err != nil {
				k.logger.Errorf("%v", err)
			}
		}
	}
	k.logger.Errorf("%v", k.connectKlcsEventStream())
}

// Subscribe implements INotificationService.
func (k *KlcsClientService) Subscribe(chBufferSize uint) async.Stream[domain.PrintJob] {
	client := async.NewBufferedStream[domain.PrintJob](chBufferSize)
	k.clients[client] = true
	return client
}

// Unsubscribe implements INotificationService.
func (k *KlcsClientService) Unsubscribe(client async.Stream[domain.PrintJob]) {
	delete(k.clients, client)
}

// -------------------------------------------------------------------------------
// Private
// -------------------------------------------------------------------------------
func (k *KlcsClientService) notify(msg async.ActionResult[domain.PrintJob]) {
	for client := range k.clients {
		client <- msg
	}
}

func (k *KlcsClientService) connectKlcsEventStream() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/printers/%s/jobs", k.klcsBackendHost, k.printerId), nil)
	if err != nil {
		return appError.NewErrNet("failed to create http request for klcs-printer-api - %v", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return appError.NewErrNet("failed to do request on klcs-printer-api - %v", err)
	}

	msg := []byte{}
LP1:
	for {
		select {
		case <-k.stop:
			break LP1
		default:
			buffer := make([]byte, READ_BUFFER_SIZE)
			n, err := resp.Body.Read(buffer)
			if err != nil {
				return appError.NewErrNet("failed to read response body from klcs-printer-api - %v", err)
			}

			data := make([]byte, n)
			copy(data, buffer)
			msg = append(msg, data...)

			// Message finished
			if n < READ_BUFFER_SIZE {
				rows := strings.Split(string(msg), "\n")
				msgType := strings.TrimPrefix(string(rows[0]), "event:")
				data := strings.TrimPrefix(string(rows[1]), "data:")

				if msgType == "message" {
					k.logger.Tracef("received job %s", data)
					job := domain.PrintJob{}
					err := json.Unmarshal([]byte(data), &job)
					if err != nil {
						return appError.NewErrDataFormat("failed to parse print job - %v", err)
					}
					k.notify(async.ActionResult[domain.PrintJob]{
						Result: job,
						Error:  nil,
					})
				}
				msg = []byte{}
			}
		}
	}
	return nil
}

// -------------------------------------------------------------------------------
// Options
// -------------------------------------------------------------------------------

// ----------------------------------------------------------------------
// Constructor
// ----------------------------------------------------------------------

func NewKlcsClientService(logger log.ILogger, klcsBackendHost string, printerId string) INotificationService[domain.PrintJob] {
	return &KlcsClientService{
		logger:          logger,
		klcsBackendHost: klcsBackendHost,
		printerId:       printerId,

		clients: map[async.Stream[domain.PrintJob]]bool{},
		stop:    make(chan any, 1),
	}
}
