package realtime

import (
	"sync"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Hub struct {
	mu      sync.RWMutex
	streams map[string]map[chan Event]struct{}
}

func NewHub() *Hub {
	return &Hub{
		streams: make(map[string]map[chan Event]struct{}),
	}
}

func (h *Hub) Subscribe(key string) chan Event {
	ch := make(chan Event, 16)

	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.streams[key]; !ok {
		h.streams[key] = make(map[chan Event]struct{})
	}
	h.streams[key][ch] = struct{}{}
	return ch
}

func (h *Hub) Unsubscribe(key string, ch chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if subs, ok := h.streams[key]; ok {
		delete(subs, ch)
		close(ch)
		if len(subs) == 0 {
			delete(h.streams, key)
		}
	}
}

func (h *Hub) Publish(key string, ev Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.streams[key] {
		select {
		case ch <- ev:
		default:
		}
	}
}
