package hw04lrucache

import "sync"

type Key string
type kvPair struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	m        *sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		m:        &sync.RWMutex{},
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.m.Lock()
	defer l.m.Unlock()
	item, ok := l.items[key]
	if ok {
		item.Value = kvPair{key: key, value: value}
		l.items[key] = item
		l.queue.MoveToFront(item)
		return true
	}

	if l.queue.Len() == l.capacity {
		last := l.queue.Back()
		l.queue.Remove(last)
		delete(l.items, last.Value.(kvPair).key)
	}

	newItem := kvPair{key: key, value: value}
	l.items[key] = l.queue.PushFront(newItem)

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.m.RLock()
	defer l.m.RUnlock()
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return item.Value.(kvPair).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.m.Lock()
	defer l.m.Unlock()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}
