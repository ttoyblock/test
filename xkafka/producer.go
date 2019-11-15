package xkafka

import (
	"math/rand"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

//var log = app.App().Log()

type Producer struct {
	Topic         string
	clientName    string
	asyncProducer sarama.SyncProducer
	partition     int32
}

func NewProducer() *Producer {
	return &Producer{}
}

func (this *Producer) Init(conf map[string]interface{}) {

	this.Topic = conf["topic"].(string)
	this.clientName = conf["clientName"].(string)
	part := conf["partition"].(int)
	this.partition = int32(part)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	broker := []string{
		conf["broker"].(string),
	}

	asyncProducer, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		println(err.Error())
	}
	this.asyncProducer = asyncProducer

	//TODO:: 没有defer
}

/**
SEND MESSAGE
*/
func (this *Producer) PutMessage(msg []byte) error {

	message := &sarama.ProducerMessage{
		Topic:     this.Topic,
		Key:       nil,
		Value:     sarama.StringEncoder(msg),
		Partition: rand.Int31n(this.partition),
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	_, _, err := this.asyncProducer.SendMessage(message)

	return err
}

func (this *Producer) Close() {
	if err := this.asyncProducer.Close(); err != nil {
		println(err.Error())
	}
}
