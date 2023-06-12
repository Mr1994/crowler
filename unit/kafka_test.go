package unit

import (
	config2 "api_client/config"
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

func TestConKafka(t *testing.T) {
	// 新建一个arama配置实例
	config := sarama.NewConfig()

	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll

	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Return.Successes = true

	// 新建一个同步生产者
	kafkaAddress := config2.IniConf.Section("KAFKA").Key("Kafka").String()
	client, err := sarama.NewSyncProducer([]string{kafkaAddress}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.Close()

	// 定义一个生产消息，包括Topic、消息内容、
	msg := &sarama.ProducerMessage{}
	msg.Topic = "k2_test3"
	//msg.Key = sarama.StringEncoder("miles")
	msg.Value = sarama.StringEncoder("hello world...")

	// 发送消息
	pid, offset, err := client.SendMessage(msg)

	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}
	fmt.Sprintf("pid:%v offset:%v\n msg:%s", pid, offset, msg.Value)
}
