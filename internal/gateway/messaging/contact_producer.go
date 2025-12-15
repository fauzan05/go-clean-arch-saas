package messaging

import (
	"golang-clean-architecture/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ContactProducer struct {
	Producer[*model.ContactEvent]
}

func NewContactProducer(producer sarama.SyncProducer, log *logrus.Logger) *ContactProducer {
	return &ContactProducer{
		Producer: Producer[*model.ContactEvent]{
			Producer: producer,
			Topic:    "contacts",
			Log:      log,
		},
	}
}
