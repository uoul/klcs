package services

import (
	"github.com/uoul/go-common/async"
	"github.com/uoul/klcs/backend/edge-printer/domain"
)

type KlcsClientService struct {
}

// Close implements INotificationService.
func (k *KlcsClientService) Close() error {
	panic("unimplemented")
}

// Run implements INotificationService.
func (k *KlcsClientService) Run() {
	panic("unimplemented")
}

// Subscribe implements INotificationService.
func (k *KlcsClientService) Subscribe(chBufferSize uint) async.Stream[domain.PrintJob] {
	panic("unimplemented")
}

// Unsubscribe implements INotificationService.
func (k *KlcsClientService) Unsubscribe(async.Stream[domain.PrintJob]) {
	panic("unimplemented")
}

func NewKlcsClientService() INotificationService[domain.PrintJob] {
	return &KlcsClientService{}
}
