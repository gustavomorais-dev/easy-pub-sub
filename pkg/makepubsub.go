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
	publisher, err := NewPublisher(ps, brokers)

	if err != nil {
		return nil, err
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.publishers = append(ps.publishers, publisher)

	return publisher, nil
}

func (ps *PubSub) CreateBroker() *Broker {
	broker := NewBroker(ps)

	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.brokers = append(ps.brokers, broker)

	return broker
}

func (bm *PubSub) GetBrokerForTopic(topic string) *Broker {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if len(bm.brokers) > 0 {
		return bm.brokers[0]
	}
	return nil
}
