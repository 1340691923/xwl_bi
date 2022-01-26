package db

import (
	"github.com/Shopify/sarama"
	"log"
	"time"
)

var KafkaSyncProducer sarama.SyncProducer
var KafkaASyncProducer sarama.AsyncProducer

var KafkaClient sarama.Client

func NewKafkaSyncProduce(host []string, username, pwd string) (conn sarama.SyncProducer, err error) {

	config := sarama.NewConfig()

	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition

	config.Version = sarama.V2_0_0_0

	if username != "" && pwd != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = pwd
		config.Net.SASL.Handshake = true
	}

	conn, err = sarama.NewSyncProducer(host, config)
	if err != nil {
		log.Println("producer closed, err:", err)
		return
	}
	return
}

func NewKafkaAsyncProduce(host []string, username, pwd string) (conn sarama.AsyncProducer, err error) {

	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.NoResponse                        // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy                  // Compress messages
	config.Producer.Flush.Frequency = time.Duration(500) * time.Millisecond // Flush batches every 500ms 不分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner               // 新选出一个partition
	config.Version = sarama.V2_0_0_0

	if username != "" && pwd != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = pwd
		config.Net.SASL.Handshake = true
	}

	conn, err = sarama.NewAsyncProducer(host, config)
	if err != nil {
		log.Println("producer closed, err:", err)
		return
	}
	return
}
