package makepubsub

import "sync"

type Broker struct {
	pubsub        *PubSub
	numReaders    int
	registrations map[string][]Subscriber
	msgCh         map[*Publisher]chan Message
	mu            sync.Mutex
}

func NewBroker(pubsub *PubSub) *Broker {
	return &Broker{
		pubsub:        pubsub,
		numReaders:    1,
		registrations: make(map[string][]Subscriber),
		msgCh:         make(map[*Publisher]chan Message),
	}
}

func (b *Broker) RegisterPublisher(publisher *Publisher) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgCh[publisher] = make(chan Message)
}

func (b *Broker) RegisterTopic(topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.registrations[topic]; !exists {
		b.registrations[topic] = []Subscriber{}
	}
}

func (b *Broker) AddSubscriberToTopic(topic string, sub *Subscriber) {
	b.RegisterTopic(topic)

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, subscriber := range b.registrations[topic] {
		if &subscriber == sub {
			return
		}
	}

	b.registrations[topic] = append(b.registrations[topic], *sub)
}

func (b *Broker) SetNumReaders(n int) {
	b.numReaders = n
}

func (b *Broker) StartReading() {
	for i := 0; i < b.numReaders; i++ {
		go b.readFromChannel()
	}
}

func (b *Broker) readFromChannel() {
	for _, channel := range b.msgCh {
		go func(ch chan Message) {
			for msg := range ch {
				b.sendMessage(msg)
			}
		}(channel)
	}
}

func (b *Broker) sendMessage(msg Message) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subscribers, ok := b.registrations[msg.topic]; ok {
		for _, subscriber := range subscribers {
			select {
			case subscriber.msgCh <- msg:
			default:
			}
		}
	}
}
