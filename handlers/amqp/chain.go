package amqp

import (
	"time"

	"github.com/Shopify/sarama"
	votingpb "github.com/egsam98/voting/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/egsam98/voting/chain/services/hyperledger"
)

type ChainHandlerOption func(*ChainHandler)

// WithTopicDead treats topic using dead letter queue pattern with provided consumption interval
func WithTopicDead(consumptionInterval time.Duration) ChainHandlerOption {
	return func(handler *ChainHandler) {
		handler.isTopicDead = true
		handler.consumptionInterval = consumptionInterval
	}
}

// ChainHandler is sarama consumer's handler to register vote into blockchain
type ChainHandler struct {
	isTopicDead         bool
	consumptionInterval time.Duration
	hyperledger         *hyperledger.Hyperledger
	producer            sarama.SyncProducer
}

func NewChainHandler(
	hyperledger *hyperledger.Hyperledger,
	producer sarama.SyncProducer,
	options ...ChainHandlerOption,
) *ChainHandler {
	h := &ChainHandler{
		hyperledger: hyperledger,
		producer:    producer,
	}

	for _, option := range options {
		option(h)
	}

	return h
}

func (c *ChainHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		vote := &votingpb.Vote{}
		if err := proto.Unmarshal(msg.Value, vote); err != nil {
			return errors.Wrapf(err, "failed to unmarshal data=%s", string(msg.Value))
		}

		log.Debug().
			Str("topic", msg.Topic).
			Int32("partition", msg.Partition).
			Int64("offset", msg.Offset).
			Interface("vote", vote).
			Msg("handlers.amqp: Received message")

		//if err := c.service.Run(session.Context(), vote); err != nil {
		//	if c.isTopicDead {
		//		time.Sleep(c.consumptionInterval)
		//		return err
		//	}
		//
		//	topicDead := msg.Topic + ".dead"
		//	if _, _, err := c.producer.SendMessage(&sarama.ProducerMessage{
		//		Topic: topicDead,
		//		Value: sarama.ByteEncoder(msg.Value),
		//	}); err != nil {
		//		return errors.Wrapf(err, "failed to send message to topic %q", topicDead)
		//	}
		//
		//	log.Error().Stack().Err(err).Msg("handlers.amqp: Vote handling error")
		//	return nil
		//}

		if c.isTopicDead {
			time.Sleep(c.consumptionInterval)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func (*ChainHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (*ChainHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
