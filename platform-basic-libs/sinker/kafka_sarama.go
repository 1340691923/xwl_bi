package sinker

import (
	"context"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/model"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

type KafkaSarama struct {
	topic     string
	cfg       model.KafkaCfg
	cg        sarama.ConsumerGroup
	sess      sarama.ConsumerGroupSession
	ctx       context.Context
	cancel    context.CancelFunc
	wgRun     sync.WaitGroup
	putFn     func(msg model.InputMessage, markFn func())
	cleanupFn func()
}

func NewKafkaSarama() *KafkaSarama {
	return &KafkaSarama{}
}

func (k *KafkaSarama) Clone() *KafkaSarama {
	return &KafkaSarama{}
}

type MyConsumerGroupHandler struct {
	k *KafkaSarama
}

func (h MyConsumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	h.k.sess = sess
	return nil
}

func (h MyConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	begin := time.Now()
	h.k.cleanupFn()
	logs.Logger.Info("consumer group cleanup",
		zap.Int32("generation id", h.k.sess.GenerationID()),
		zap.Duration("cost", time.Since(begin)))
	return nil
}

func (h MyConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.k.putFn(model.InputMessage{
			Topic:     msg.Topic,
			Partition: int(msg.Partition),
			Key:       msg.Key,
			Value:     msg.Value,
			Offset:    msg.Offset,
			Timestamp: &msg.Timestamp,
		}, func() {
			sess.MarkMessage(msg, "")
		})

	}
	return nil
}

func (k *KafkaSarama) Init(cfg model.KafkaCfg, topicName, consumerGroup string, putFn func(msg model.InputMessage, markFn func()), cleanupFn func()) (err error) {
	k.cfg = cfg
	k.ctx, k.cancel = context.WithCancel(context.Background())
	k.putFn = putFn
	k.cleanupFn = cleanupFn
	k.topic = topicName
	sarCfg, err := GetSaramaConfig(cfg)
	if err != nil {
		return err
	}
	sarCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	cg, err := sarama.NewConsumerGroup(cfg.Addresses, consumerGroup, sarCfg)
	if err != nil {
		return err
	}

	k.cg = cg
	return nil
}

func GetSaramaConfig(kfkCfg model.KafkaCfg) (sarCfg *sarama.Config, err error) {
	sarCfg = sarama.NewConfig()
	sarCfg.Version = sarama.V2_0_0_0
	sarCfg.Consumer.Return.Errors = false
	// check for authentication
	if kfkCfg.Username != "" && kfkCfg.Password != "" {
		sarCfg.Net.SASL.Enable = true
		sarCfg.Net.SASL.User = kfkCfg.Username
		sarCfg.Net.SASL.Password = kfkCfg.Password
	}
	sarCfg.ChannelBufferSize = 1024
	return
}

func (k *KafkaSarama) Run() {
	k.wgRun.Add(1)
	defer k.wgRun.Done()
LOOP_SARAMA:
	for {

		handler := MyConsumerGroupHandler{k}

		if k.ctx.Err() != nil {
			return
		}

		if err := k.cg.Consume(k.ctx, []string{k.topic}, handler); err != nil {
			if errors.Is(err, context.Canceled) {
				logs.Logger.Info("KafkaSarama.Run quit due to context has been canceled", zap.String("task", k.topic))
				break LOOP_SARAMA
			} else if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				logs.Logger.Info("KafkaSarama.Run quit due to consumer group has been closed", zap.String("task", k.topic))
				break LOOP_SARAMA
			} else {
				logs.Logger.Error("sarama.ConsumerGroup.Consume failed", zap.String("task", k.topic), zap.Error(err))
				continue
			}
		}
	}
}

func (k *KafkaSarama) CommitMessages(msg *model.InputMessage) error {
	k.sess.MarkOffset(msg.Topic, int32(msg.Partition), msg.Offset+1, "")
	return nil
}

func (k *KafkaSarama) Stop() error {
	k.cancel()
	k.cg.Close()
	k.wgRun.Wait()
	return nil
}

func (k *KafkaSarama) Description() string {
	return "kafka consumer of topic " + k.topic
}
