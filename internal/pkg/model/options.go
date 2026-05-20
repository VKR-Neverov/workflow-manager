package model

import "sync"

type Options interface {
	SetString(key, value string)
	GetString(key string) string
}

type options struct {
	mu      sync.RWMutex
	options map[string]string
}

func NewOptions() *options {
	return &options{
		options: make(map[string]string),
	}
}

func (o *options) SetString(key, value string) {
	o.mu.Lock()
	o.options[key] = value
	o.mu.Unlock()
}

func (o *options) GetString(key string) string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.options[key]
}
