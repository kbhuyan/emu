package util

import "sync"

type PubSub[S comparable, T any] struct {
	mu          sync.Mutex
	subscribers map[S]map[chan T]bool
}

// NewPubSub initializes a new PubSub instance with a map to hold subscribers for each topic.
// The map key is of type S (the topic type), and the value is a map of channels of type T.
// Each channel represents a subscriber for that topic.
func NewPubSub[S comparable, T any]() *PubSub[S, T] {
	return &PubSub[S, T]{
		subscribers: make(map[S]map[chan T]bool),
	}
}

// Subscribe creates a new channel for the given topic and returns it.
func (ps *PubSub[S, T]) Subscribe(topic S) chan T {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.subscribers[topic]; !ok {
		ps.subscribers[topic] = make(map[chan T]bool)
	}
	ch := make(chan T, 1)
	ps.subscribers[topic][ch] = true
	return ch
}

// Close unsubscribes the channel from the topic and cleans up if no more subscribers exist for that topic.

func (ps *PubSub[S, T]) Close(topic S, ch <-chan T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	// Remove the channel from the subscribers map
	if subs, ok := ps.subscribers[topic]; ok {
		for s := range subs {
			if s == ch {
				delete(subs, s)
				break
			}
		}
		if len(subs) == 0 {
			delete(ps.subscribers, topic)
		}
	}
}

func (ps *PubSub[S, T]) Publish(topic S, val T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if subs, ok := ps.subscribers[topic]; ok {
		for s := range subs {
			s <- val
		}
	}
}
