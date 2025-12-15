package config

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewKafkaConsumerGroup(config *viper.Viper, log *logrus.Logger) sarama.ConsumerGroup {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	offsetReset := config.GetString("kafka.auto.offset.reset")
	if offsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	brokers := strings.Split(config.GetString("kafka.bootstrap.servers"), ",")
	groupID := config.GetString("kafka.group.id")

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	return consumerGroup
}

func NewKafkaProducer(config *viper.Viper, log *logrus.Logger) sarama.SyncProducer {
	if !config.GetBool("kafka.producer.enabled") {
		log.Info("Kafka producer is disabled")
		return nil
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3

	brokers := strings.Split(config.GetString("kafka.bootstrap.servers"), ",")

	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	return producer
}
