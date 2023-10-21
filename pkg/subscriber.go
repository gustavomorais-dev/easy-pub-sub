package makepubsub

import "sync"

type Subscriber struct {
	topics []string
	msgCh  chan Message
	mu     sync.Mutex
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		msgCh:  make(chan Message),
		topics: []string{},
	}
}

func (s *Subscriber) SubscribeTo(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.topics = append(s.topics, topic)
	// lógica para registrar o tópico no broker
}

func (s *Subscriber) UnsubscribeFrom(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove o tópico da lista de tópicos do subscriber
	for i, t := range s.topics {
		if t == topic {
			s.topics = append(s.topics[:i], s.topics[i+1:]...)
			break
		}
	}
	// lógica para cancelar o registro do tópico no broker
}
