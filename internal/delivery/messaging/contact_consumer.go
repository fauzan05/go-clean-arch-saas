package messaging

import (
	"encoding/json"
	"golang-clean-architecture/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ContactConsumer struct {
	Log *logrus.Logger
}

func NewContactConsumer(log *logrus.Logger) *ContactConsumer {
	return &ContactConsumer{
		Log: log,
	}
}

func (c ContactConsumer) Consume(message *sarama.ConsumerMessage) error {
	ContactEvent := new(model.ContactEvent)
	if err := json.Unmarshal(message.Value, ContactEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling Contact event")
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic contacts with event: %v from partition %d", ContactEvent, message.Partition)
	return nil
}
