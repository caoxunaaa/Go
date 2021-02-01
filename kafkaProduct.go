package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 基于sarama第三方库开发的kafka client

func main() {
	// 新建一个arama配置实例
	config := sarama.NewConfig()

	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll

	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Return.Successes = true

	// 新建一个同步生产者
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.Close()

	// 定义一个生产消息，包括Topic、消息内容、
	msg := &sarama.ProducerMessage{}
	msg.Topic = "revolution"
	msg.Key = sarama.StringEncoder("miles")
	msg.Value = sarama.StringEncoder("123hello world...")

	// 发送消息
	pid, offset, err := client.SendMessage(msg)

	msg2 := &sarama.ProducerMessage{}
	msg2.Topic = "revolution"
	msg2.Key = sarama.StringEncoder("monroe")
	msg2.Value = sarama.StringEncoder("321hello world2...")
	pid2, offset2, err := client.SendMessage(msg2)

	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
	fmt.Printf("pid2:%v offset2:%v\n", pid2, offset2)
}
