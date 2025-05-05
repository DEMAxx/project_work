package lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*cacheItem
}

type cacheItem struct {
	key   Key
	value interface{}
	item  *ListItem
}

var mutex sync.Mutex

func (lruCache *lruCache) Set(key Key, value interface{}) bool {
	mutex.Lock()
	defer mutex.Unlock()

	item, ok := lruCache.items[key]

	if ok {
		item.value = value
		lruCache.queue.MoveToFront(item.item)

		return true
	}

	if lruCache.capacity == lruCache.queue.Len() {
		back := lruCache.queue.Back()
		valKey, ok := back.Value.(Key)

		if !ok {
			return false
		}
		delete(lruCache.items, valKey)
		lruCache.queue.Remove(back)
	}
	newItem := lruCache.queue.PushFront(key)

	lruCache.items[key] = &cacheItem{
		key:   key,
		value: value,
		item:  newItem,
	}

	return false
}

func (lruCache *lruCache) Get(key Key) (interface{}, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	item, ok := lruCache.items[key]

	if !ok {
		return nil, false
	}

	lruCache.queue.MoveToFront(item.item)
	return item.value, true
}

func (lruCache *lruCache) Clear() {
	mutex.Lock()
	defer mutex.Unlock()

	lruCache.queue = new(list)
	lruCache.items = make(map[Key]*cacheItem, lruCache.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    new(list),
		items:    make(map[Key]*cacheItem, capacity),
	}
}
