package messaging

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"task_tracker/services"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

type UserInfoMsg struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserInfoConsumer struct {
	userInfoChannel chan services.UserInfo
	brokersUrl      []string
	topic           string
}

func NewUserInfoConsumer() UserInfoConsumer {
	return UserInfoConsumer{
		userInfoChannel: make(chan services.UserInfo),
	}
}

func (uic *UserInfoConsumer) GetUserInfoChannel() <-chan services.UserInfo {
	return uic.userInfoChannel
}

func (uic *UserInfoConsumer) ReadFromTopic(brokersUrl []string, topic string) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	fmt.Printf("Connecting to brokers at url: %#v\n", brokersUrl)
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return err
	}
	consumer, err := conn.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	fmt.Println("Consumer started and reading from topic: " + topic)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				fmt.Printf("Received message: %s \n", string(msg.Value))
				userInfoMsg := UserInfoMsg{}
				parseErr := json.Unmarshal(msg.Value, &userInfoMsg)
				if parseErr != nil {
					fmt.Println("Parsing error is detected: " + parseErr.Error())
				} else {
					uic.userInfoChannel <- services.UserInfo{
						Id:   userInfoMsg.Id,
						Name: userInfoMsg.Name,
					}
				}
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
				return
			}
		}
	}()

	<-doneCh

	fmt.Println("Stopping consumer")
	return conn.Close()
}
