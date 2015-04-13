package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Subscriber interface {
	Notify(messgae interface{})
}

type Publisher interface {
	Subscribe(key string, subscriber Subscriber)
	Publish(key string, message interface{})
}

type PuhblishDelegate func(message interface{})

type DefaultSubscriber struct {
	channel chan interface{}
}

func NewSubscriber(publishDelegate PuhblishDelegate) (subscriber Subscriber) {
	defaultSubscriber := &DefaultSubscriber{make(chan interface{})}
	go defaultSubscriber.handlePublications(publishDelegate)
	return defaultSubscriber
}

func (subscriber *DefaultSubscriber) handlePublications(publishDelegate PuhblishDelegate) {
	for {
		select {
		case message, ok := <-subscriber.channel:
			if ok {
				publishDelegate(message)
			} else {
				subscriber.channel = nil
				return
			}
		default:
			time.Sleep(1)
		}
	}
}

func (subscriber *DefaultSubscriber) Notify(message interface{}) {
	subscriber.channel <- message
}

type DefaultPublisher struct {
	subscribers map[string][]Subscriber
}

func (publisher *DefaultPublisher) Subscribe(key string, subscriber Subscriber) {
	if subscriberList := publisher.subscribers[key]; subscriberList == nil {
		fmt.Printf("initializing the subscriber list for key: %s\n", key)
		publisher.subscribers[key] = []Subscriber{}
	}

	fmt.Println("Adding subcriber")
	publisher.subscribers[key] = append(publisher.subscribers[key], subscriber)
	fmt.Printf("Subscriber count = %d\n", len(publisher.subscribers[key]))
}

func (publisher *DefaultPublisher) Publish(key string, message interface{}) {
	subscriberList := publisher.subscribers[key]

	fmt.Printf("Notifying %d subscribers\n", len(subscriberList))
	for _, subscriber := range subscriberList {
		if subscriber != nil {
			subscriber.Notify(message)
		}
	}
}

func NewPublisher() (publisher Publisher) {
	publisher = &DefaultPublisher{make(map[string][]Subscriber)}
	return
}

func Test_PublishingAMessageToMultipleSubscribers(t *testing.T) {

	var wg sync.WaitGroup

	publisher := NewPublisher()

	wg.Add(1)
	subscriber1 := NewSubscriber(func(value interface{}) {
		fmt.Printf("Subscriber 1 received %v\n", value)
		wg.Done()
	})

	wg.Add(1)
	subscriber2 := NewSubscriber(func(value interface{}) {
		fmt.Printf("Subscriber 2 received %v\n", value)
		wg.Done()
	})

	publisher.Subscribe("foobar", subscriber1)
	publisher.Subscribe("foobar", subscriber2)
	publisher.Publish("foobar", 1)

	wg.Wait()
}
