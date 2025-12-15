package messaging

import (
	"encoding/json"
	"golang-clean-architecture/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type UserConsumer struct {
	Log *logrus.Logger
}

func NewUserConsumer(log *logrus.Logger) *UserConsumer {
	return &UserConsumer{
		Log: log,
	}
}

func (c UserConsumer) Consume(message *sarama.ConsumerMessage) error {
	UserEvent := new(model.UserEvent)
	if err := json.Unmarshal(message.Value, UserEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling User event")
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic users with event: %v from partition %d", UserEvent, message.Partition)
	return nil
}
