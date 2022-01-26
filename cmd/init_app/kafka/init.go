package kafka

import (
	"fmt"
	"github.com/1340691923/xwl_bi/model"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

//初始化kafka数据
func Init() {
	config := sarama.NewConfig()

	config.Version = sarama.V2_0_0_0
	if model.GlobConfig.Comm.Kafka.Username != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = model.GlobConfig.Comm.Kafka.Username
		config.Net.SASL.Password = model.GlobConfig.Comm.Kafka.Password
		config.Net.SASL.Handshake = true
	}

	config.Consumer.Group.Session.Timeout = 15 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 5 * time.Second

	conn, err := sarama.NewClusterAdmin(model.GlobConfig.Comm.Kafka.Addresses, config)
	if err != nil {
		log.Println(fmt.Sprintf("kafka 链接初始化失败:%s", err.Error()))
		panic(err)
	}
	s, err := conn.ListTopics()
	for topic := range s {
		log.Println("您所拥有的TOPIC为：", topic)
	}

	if _, ok := s[model.GlobConfig.Comm.Kafka.ReportTopicName]; !ok {
		detail := sarama.TopicDetail{NumPartitions: model.GlobConfig.Comm.Kafka.NumPartitions, ReplicationFactor: 1}
		err = conn.CreateTopic(model.GlobConfig.Comm.Kafka.ReportTopicName, &detail, false)
		if err != nil {
			log.Println("创建TOPIC失败！", model.GlobConfig.Comm.Kafka.ReportTopicName)
			panic(err)
		}

		err = conn.Close()
		if err != nil {
			panic(err)
		}
		log.Println("初始化TOPIC完成！", model.GlobConfig.Comm.Kafka.ReportTopicName)
	} else {
		log.Println("您已拥有该TOPIC：", model.GlobConfig.Comm.Kafka.ReportTopicName)
	}
}
