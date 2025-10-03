package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type Comment struct {
	Message string `json:"message"`
}

func main() {
	r := gin.Default()
	commentgroup := r.Group("/comments")
	commentgroup.POST("/", handleComment)
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func handleComment(c *gin.Context) {
	topic := "comment"
	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(500, gin.H{"error": "Unable to read payload"})
	}

	err := SendMessage(topic, comment)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message:": "Send message to kafka successfully"})
}

func SendMessage(topic string, message any) error {
	producer, err := ConnectToBroker()
	if err != nil {
		return err
	}

	stringfyMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(stringfyMessage),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message stored at topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func ConnectToBroker() (sarama.SyncProducer, error) {
	brokersUrls := []string{"localhost:29092"}
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	conn, err := sarama.NewSyncProducer(brokersUrls, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
