package makepubsub

import (
	"errors"
	"sync"
)

type Publisher struct {
	brokers []Broker
	mu      sync.Mutex
}

func NewPublisher(brokers []Broker) (*Publisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("no brokers provided")
	}

	publisher := Publisher{
		brokers: brokers,
	}

	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	for _, broker := range brokers {
		broker.RegisterPublisher(&publisher)
	}

	return &publisher, nil
}

func (p *Publisher) Publish(msg Message) {
	p.mu.Lock()
	defer p.mu.Unlock()

	updatedBrokers := make(map[*Broker]struct{})

	for _, broker := range p.brokers {
		broker.mu.Lock()

		if _, ok := updatedBrokers[&broker]; !ok {
			updatedBrokers[&broker] = struct{}{}
			for key := range broker.msgCh {
				broker.msgCh[key] <- msg
			}
		}

		broker.mu.Unlock()
	}
}
