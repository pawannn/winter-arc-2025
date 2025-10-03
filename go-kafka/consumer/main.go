package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	topic := "comment"
	consumer, err := connectToConsumer()
	if err != nil {
		log.Fatal(err)
	}

	partition, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}
	defer partition.Close()

	doneChan := make(chan os.Signal, 1)
	signal.Notify(doneChan, syscall.SIGTERM, syscall.SIGINT)
	defer close(doneChan)

	messageCount := 0

	go func() {
		for {
			select {
			case <-doneChan:
				return
			case message := <-partition.Messages():
				msg := message.Value
				messageCount += 1
				fmt.Println("received message from kafka : ", string(msg))
			}
		}
	}()

	fmt.Println("Comments consumer started...")
	<-doneChan
	fmt.Println("Total message received : ", messageCount)
}

func connectToConsumer() (sarama.Consumer, error) {
	brokersUrls := []string{"localhost:29092"}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokersUrls, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
