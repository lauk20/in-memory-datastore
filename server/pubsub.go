package main

import (
	"context"
	"sync"

	pubsubpb "datastore/protos/pubsub"
)

// pub/sub message broker
type Broker struct {
	mutex       sync.RWMutex             // mutex protecting subscribers
	subscribers map[string][]chan string // subscribers to a specific topic, each with own channel
}

// create a new Broker
// returns a new *Broker
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan string),
	}
}

// Server for the pubsub service
type PubSubServer struct {
	pubsubpb.UnimplementedPubSubServer

	// message broker
	broker *Broker
}

// Publish to pubsub server
func (s *PubSubServer) Publish(ctx context.Context, message *pubsubpb.Pub) (*pubsubpb.NumSubs, error) {
	s.broker.mutex.RLock()
	defer s.broker.mutex.RUnlock()

	var count int32
	count = 0
	for _, channel := range s.broker.subscribers[message.Topic] {
		channel <- message.Msg
		count += 1
	}

	return &pubsubpb.NumSubs{Value: count}, nil
}

// Subscribe to pubsub server
func (s *PubSubServer) Subscribe(topic *pubsubpb.String, stream pubsubpb.PubSub_SubscribeServer) error {
	s.broker.mutex.RLock()
	channel := make(chan string)
	s.broker.subscribers[topic.Msg] = append(s.broker.subscribers[topic.Msg], channel)
	s.broker.mutex.RUnlock()

	for {
		result := <-channel
		if err := stream.Send(&pubsubpb.String{Msg: result}); err != nil {
			return err
		}
	}

	return nil
}
