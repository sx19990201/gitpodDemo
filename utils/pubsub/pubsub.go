package pubsub

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"sync"
)

const (
	TopicName            = "Reload"
	TopicPayloadStarting = "Starting"
	TopicPayloadBuilding = "Building"
	TopicPayloadStop     = "Stop"
)

var GoChannel *gochannel.GoChannel
var once sync.Once

var MessageChan <-chan *message.Message
var messageChanOnce sync.Once

func GetMessageChan(message <-chan *message.Message) <-chan *message.Message {
	if MessageChan == nil {
		messageChanOnce.Do(func() {
			MessageChan = message
		})
	}
	return MessageChan
}

func GetPubSub() *gochannel.GoChannel {
	if GoChannel == nil {
		once.Do(func() {
			GoChannel = gochannel.NewGoChannel(
				gochannel.Config{},
				watermill.NewStdLogger(false, false),
			)
		})
	}
	return GoChannel
}

func GetTopicMessage(pubSub *gochannel.GoChannel, topicName string) <-chan *message.Message {
	messages, err := pubSub.Subscribe(context.Background(), topicName)
	if err != nil {
		panic(err)
	}
	return messages
}

func PublishMessages(topicName, payload string) {
	msg := message.NewMessage(watermill.NewUUID(), []byte(payload))
	if err := GetPubSub().Publish(topicName, msg); err != nil {
		panic(err)
	}
}
