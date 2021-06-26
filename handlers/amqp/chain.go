package amqp

import (
	"time"

	"github.com/Shopify/sarama"
	votingpb "github.com/egsam98/voting/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/egsam98/voting/chain/services/smart"
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
	client              *smart.Client
	producer            sarama.SyncProducer
}

func NewChainHandler(
	client *smart.Client,
	producer sarama.SyncProducer,
	options ...ChainHandlerOption,
) *ChainHandler {
	h := &ChainHandler{
		client:   client,
		producer: producer,
	}

	for _, option := range options {
		option(h)
	}

	return h
}

func (c *ChainHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		select {
		case <-session.Context().Done():
			return nil
		default:
		}

		if err := c.processMessage(msg); err != nil {
			return err
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *ChainHandler) processMessage(msg *sarama.ConsumerMessage) error {
	log.Debug().
		Str("topic", msg.Topic).
		Int32("partition", msg.Partition).
		Int64("offset", msg.Offset).
		Msg("amqp: Received message")

	vote, err := c.decodeToVote(msg)
	if err != nil {
		log.Error().
			Stack().
			Err(err).
			Msg("amqp: Invalid vote")
		return nil
	}

	if err := c.client.RegisterVote(vote); err != nil {
		if errors.Is(err, smart.ErrInvalidVote) {
			log.Warn().
				Err(err).
				Msg("amqp: Invalid vote")
			return nil
		}

		log.Error().Stack().Err(err).Msg("amqp: Vote handling error")

		if c.isTopicDead {
			time.Sleep(c.consumptionInterval)
			return err
		}

		topicDead := msg.Topic + ".dead"
		if _, _, err := c.producer.SendMessage(&sarama.ProducerMessage{
			Topic: topicDead,
			Value: sarama.ByteEncoder(msg.Value),
		}); err != nil {
			return errors.Wrapf(err, "failed to send message to topic %q", topicDead)
		}

		return nil
	}

	log.Info().Msg("amqp: Vote has been registered")

	if c.isTopicDead {
		time.Sleep(c.consumptionInterval)
	}

	return nil
}

func (*ChainHandler) decodeToVote(msg *sarama.ConsumerMessage) (*votingpb.Vote, error) {
	vote := &votingpb.Vote{}
	if err := proto.Unmarshal(msg.Value, vote); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal body to %T", vote)
	}
	if vote.CandidateId == 0 {
		return nil, errors.New("vote has no candidate ID")
	}
	if vote.Voter == nil {
		return nil, errors.New("vote has no voter")
	}
	if vote.Voter.Passport == "" {
		return nil, errors.New("voter has no passport")
	}
	if vote.Voter.Fullname == "" {
		return nil, errors.New("voter has no fullname")
	}
	if vote.Voter.BirthDate == 0 {
		return nil, errors.New("voter has no birth date")
	}
	return vote, nil
}

func (*ChainHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (*ChainHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
