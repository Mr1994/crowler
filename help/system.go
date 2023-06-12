package help

import (
	"api_client/config"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
	"os"
	"runtime"
)

func SetAppIniFile(source string) *ini.File {

	// 获取gopath
	path := GetAppIni(source)
	flag.StringVar(&Appini, "c", path, "configure file")
	flag.Parse()
	var err error
	config.IniConf, err = ini.Load(path) //此路径为ini文件的路径
	if err != nil {
		fmt.Println("读取ini文件失败", err)
	}
	return config.IniConf
}

// 同步发送kafka消息
func (c *Common) PushSyncKafkaMessage(topic string, message interface{}) {
	client := config.KafkaSyncClient
	//defer client.Close()

	// 定义一个生产消息，包括Topic、消息内容、
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	//msg.Key = sarama.StringEncoder("miles")
	messageString := c.InterFaceToString(message)
	msg.Value = sarama.StringEncoder(messageString)

	// 发送消息
	pid, offset, _ := client.SendMessage(msg)
	str := fmt.Sprintf("pid:%v offset:%v\n msg:%s", pid, offset, msg.Value)
	c.Log(str, "kafkaInfo")
}

// 异步发送kafka消息
func (c *Common) PushAsyncKafkaMessage(topic string, message interface{}) {
	producer := config.KafkaAsyncClient
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case <-p.Successes():
				//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
			}
		}
	}(producer)
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	messageString := c.InterFaceToString(message)
	msg.Value = sarama.StringEncoder(messageString)
	//使用通道发送
	producer.Input() <- msg

}

// ClusterConsumer
//  @Description:
//  @receiver c
//  @param topics
//  @param groupId
func (c *Common) ClusterConsumer(topics string, groupId string) {

	//defer consumer.Close()
	partitionConsumer, err := config.KafkaConsumerClient.ConsumePartition(topics, 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("try create partitionConsumer error %s\n", err.Error())
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
				msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Printf("err :%s\n", err.Error())
		}
	}

}

// GetAppIni
//  @Description: 获取ini基本配置
//  @param source  数据来源
//  @return string
func GetAppIni(source string) string {
	helpCommon := new(Common)
	pwd, _ := os.Getwd()
	fmt.Println(pwd)

	dir := helpCommon.GetParentDirectory(pwd)
	sysType := runtime.GOOS
	fmt.Println(dir + "/config/app.ini")

	if sysType == "windows" {
		if source == "main" {
			dir = dir + "\\goTest"
		}
		return dir + "\\config\\app.ini"
	} else {
		if source == "main" {
			dir = dir + "/goApi/"
		}
		return dir + "/config/app.ini"
	}

}
