package makepubsub

import (
	"sync"
)

type PubSub struct {
	publishers  []*Publisher
	brokers     []*Broker
	subscribers []*Subscriber
	mu          sync.Mutex
}

func Start() *PubSub {
	return &PubSub{
		publishers:  []*Publisher{},
		brokers:     []*Broker{},
		subscribers: []*Subscriber{},
	}
}

func (ps *PubSub) CreatePublisher(brokers []Broker) (*Publisher, error) {
	publisher, err := NewPublisher(brokers)

	if err != nil {
		return nil, err
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.publishers = append(ps.publishers, publisher)

	return publisher, nil
}

func (ps *PubSub) CreateBroker() *Broker {
	broker := NewBroker()

	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.brokers = append(ps.brokers, broker)

	return broker
}
