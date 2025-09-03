package kafka

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/shopify/sarama"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	logger   *logrus.Logger
	handler  ConsumerHandler
}

type ConsumerHandler interface {
	HandleTransferEvent(event TransferEvent) error
}

func NewConsumer(groupID string, handler ConsumerHandler, logger *logrus.Logger) (*Consumer, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		logger:   logger,
		handler:  handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	topics := []string{"transfer-events"}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := c.consumer.Consume(ctx, topics, c); err != nil {
				c.logger.WithError(err).Error("Error consuming messages")
				return err
			}
		}
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			var event TransferEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				c.logger.WithError(err).Error("Error unmarshaling transfer event")
				session.MarkMessage(message, "")
				continue
			}

			if err := c.handler.HandleTransferEvent(event); err != nil {
				c.logger.WithError(err).Error("Error handling transfer event")
			} else {
				c.logger.WithField("requestId", event.RequestID).Info("Transfer event processed successfully")
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
