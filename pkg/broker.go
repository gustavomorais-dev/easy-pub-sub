package makepubsub

import "sync"

type Broker struct {
	numReaders    int
	registrations map[string][]Subscriber
	msgCh         map[*Publisher]chan Message
	mu            sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
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
