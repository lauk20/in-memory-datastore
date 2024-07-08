package main

import (
	"sync"
)

// pub/sub message broker
type Broker struct {
	mutex       sync.RWMutex            // mutex protecting subscribers
	subscribers map[string][]chan Value // subscribers to a specific topic, each with own channel
}

// create a new Broker
// returns a new *Broker
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan Value),
	}
}

// publish message v Value to all channels under a single topic
func (broker *Broker) Publish(topic string, v Value) {
	broker.mutex.RLock()
	defer broker.mutex.RUnlock()

	for _, channel := range broker.subscribers[topic] {
		channel <- v
	}
}

// create new subscriber to a topic
// creates a channel for the subscriber
// returns the channel created
func (broker *Broker) Subscribe(topic string) chan Value {
	broker.mutex.RLock()
	defer broker.mutex.RUnlock()

	channel := make(chan Value)
	broker.subscribers[topic] = append(broker.subscribers[topic], channel)
	return channel
}
