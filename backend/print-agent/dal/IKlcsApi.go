package dal

import "context"

type IKlcsApi interface {
	AcknowledgePrintJob(ctx context.Context, printerId string, transactionId string) error
}
