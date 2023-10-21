package makepubsub

import (
	"errors"
	"sync"
)

type Publisher struct {
	brokers map[*Broker]chan Message
	mu      sync.Mutex
}

func NewPublisher(brokers []*Broker) (*Publisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("no brokers provided")
	}

	brokerMap := make(map[*Broker]chan Message)
	for _, broker := range brokers {
		brokerMap[broker] = make(chan Message)
	}

	return &Publisher{
		brokers: brokerMap,
	}, nil
}

func (p *Publisher) Publish(msg Message) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, ch := range p.brokers {
		ch <- msg
	}
}
