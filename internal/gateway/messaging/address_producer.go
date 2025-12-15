package messaging

import (
	"golang-clean-architecture/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type AddressProducer struct {
	Producer[*model.AddressEvent]
}

func NewAddressProducer(producer sarama.SyncProducer, log *logrus.Logger) *AddressProducer {
	return &AddressProducer{
		Producer: Producer[*model.AddressEvent]{
			Producer: producer,
			Topic:    "addresses",
			Log:      log,
		},
	}
}
