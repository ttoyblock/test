package xkafka

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

type Consumer struct {
	Topic     string
	GroupName string
	consumer  *cluster.Consumer
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (this *Consumer) Init(conf map[string]interface{}) {
	this.Topic = conf["topic"].(string)
	this.GroupName = conf["groupName"].(string)

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //new?
	//config.Net.MaxOpenRequests = 5
	//config.Net.DialTimeout = 30*1000
	//config.Net.ReadTimeout = 30*1000
	//config.Net.WriteTimeout = 30*1000

	brokers := []string{conf["broker"].(string)}
	topics := []string{this.Topic}
	consumer, err := cluster.NewConsumer(brokers, this.GroupName, topics, config)
	if err != nil {
		panic(err.Error())
	}

	this.consumer = consumer
	//defer consumer.Close()

}

func (this *Consumer) RunByCall(callBack func(message *sarama.ConsumerMessage) error) {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for err := range this.consumer.Errors() {
			println(err.Error())
		}
	}()

	go func() {
		for ntf := range this.consumer.Notifications() {

			fmt.Println("kafka rebalance info" + ntf.Type.String())

		}
	}()

	for {
		select {
		case msg, ok := <-this.consumer.Messages():
			{
				if ok {

					err := callBack(msg)
					if err == nil {
						this.consumer.MarkOffset(msg, "")
					} else {
						fmt.Println(string(msg.Value))
					}

				} else {
					println("not Message")
				}
			}
		case <-signals:
			return
		}
	}

}

func (this *Consumer) Close() {
	this.consumer.Close()
}
