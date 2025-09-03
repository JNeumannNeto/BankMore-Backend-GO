package kafka

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/shopify/sarama"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   *logrus.Logger
}

type TransferEvent struct {
	RequestID               string  `json:"requestId"`
	OriginAccountID         string  `json:"originAccountId"`
	DestinationAccountID    string  `json:"destinationAccountId"`
	DestinationAccountNumber string  `json:"destinationAccountNumber"`
	Amount                  float64 `json:"amount"`
	TransferID              string  `json:"transferId"`
}

func NewProducer(logger *logrus.Logger) (*Producer, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		logger:   logger,
	}, nil
}

func (p *Producer) PublishTransferEvent(event TransferEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		p.logger.WithError(err).Error("Failed to marshal transfer event")
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "transfer-events",
		Value: sarama.StringEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.WithError(err).Error("Failed to send transfer event to Kafka")
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"partition": partition,
		"offset":    offset,
		"requestId": event.RequestID,
	}).Info("Transfer event published successfully")

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
