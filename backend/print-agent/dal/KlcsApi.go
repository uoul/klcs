package dal

import (
	"context"
	"fmt"
	"net/http"

	appError "github.com/uoul/klcs/backend/print-agent/error"
)

type KlcsApi struct {
	klcsBackendHost string
	httpClient      *http.Client
}

// AcknowledgePrintJob implements IKlcsApi.
func (k *KlcsApi) AcknowledgePrintJob(ctx context.Context, printerId string, transactionId string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/", k.klcsBackendHost), nil)
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

func NewKlcsApi(klcsBackendHost string) IKlcsApi {
	return &KlcsApi{
		klcsBackendHost: klcsBackendHost,
	}
}
