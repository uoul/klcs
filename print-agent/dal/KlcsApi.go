package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/uoul/go-async"
	"github.com/uoul/klcs/backend/print-agent/domain"
	appError "github.com/uoul/klcs/backend/print-agent/error"
)

type KlcsApi struct {
	klcsBackendHost string
	httpClient      *http.Client

	readBufferSize int
}

func (k *KlcsApi) GetPrintJobs(ctx context.Context, printerId string) (async.Sequence[domain.PrintJob], error) {
	// Create HTTP Request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/printers/%s/jobs", k.klcsBackendHost, printerId), nil)
	if err != nil {
		return nil, appError.NewErrNet("failed to create http request for klcs-printer-api - %v", err)
	}
	// Set HTTP Headers
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	// Do HTTP Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, appError.NewErrNet("failed to do request on klcs-printer-api - %v", err)
	}
	// Listen on incomming messages
	slog.Info("Listening for printJobs...", slog.Any("url", resp.Request.URL))
	jobs := make(async.Sequence[domain.PrintJob])
	go func() {
		msg := []byte{}
		buffer := make([]byte, k.readBufferSize)
		for {
			// Read from connection
			n, err := resp.Body.Read(buffer)
			if err != nil {
				jobs <- async.Fail[domain.PrintJob](appError.NewErrNet("failed to read response body from klcs-printer-api - %v", err))
				close(jobs)
				return
			}
			data := make([]byte, n)
			copy(data, buffer)
			msg = append(msg, data...)
			// Message finished (if buffer is not used completely)
			if n < k.readBufferSize {
				rows := strings.Split(string(msg), "\n")
				msgType := strings.TrimPrefix(string(rows[0]), "event:")
				data := strings.TrimPrefix(string(rows[1]), "data:")
				if msgType == "message" {
					slog.Debug("New incomming printjob", slog.String("job", string(data)))
					job := domain.PrintJob{}
					err := json.Unmarshal([]byte(data), &job)
					if err != nil {
						jobs <- async.Fail[domain.PrintJob](appError.NewErrDataFormat("failed to parse print job - %v", err))
						close(jobs)
						return
					}
					jobs <- async.Success(job)
				}
				msg = []byte{}
			}
		}
	}()
	return jobs, nil
}

// AcknowledgePrintJob send acknowledement request to klcs api
func (k *KlcsApi) AcknowledgePrintJob(ctx context.Context, printerId string, transactionId string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/printers/%s/jobs/acknowledgement/%s", k.klcsBackendHost, printerId, transactionId), nil)
	if err != nil {
		return appError.NewErrNet("failed to create http request - %v", err)
	}
	resp, err := k.httpClient.Do(req)
	if err != nil {
		return appError.NewErrNet("failed to acknowledge printer job - %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return appError.NewErrNet("failed to acknowledge printer job - StatusCode: %d", resp.StatusCode)
	}
	return nil
}

func NewKlcsApi(klcsBackendHost string) *KlcsApi {
	return &KlcsApi{
		klcsBackendHost: klcsBackendHost,
		httpClient:      http.DefaultClient,

		readBufferSize: 1024,
	}
}
