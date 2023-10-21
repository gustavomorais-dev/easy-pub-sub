package makepubsub

import "sync"

type Broker struct {
	subscribers map[string][]chan string
	mu          sync.Mutex
}
