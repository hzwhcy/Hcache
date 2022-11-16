package lru

import "container/list"

type Cache struct {
	capacity  int64
	size      int64
	l         *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}
type Value interface {
	Len() int
}
type entry struct {
	key   string
	value Value
}

// New 新建一个缓存
func New(capacity int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		capacity:  capacity,
		l:         list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Add 增加一个元素
func (c *Cache) Add(key string, value Value) {
	if val, ok := c.cache[key]; ok {
		c.l.MoveToFront(val)
		entryVal := val.Value.(*entry)
		c.size += int64(value.Len()) - int64(entryVal.value.Len())
		entryVal.value = value
	} else {
		ele := c.l.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.size += int64(len(key)) + int64(value.Len())
	}
	for c.capacity != 0 && c.capacity < c.size {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.l.Back()
	if ele != nil {
		c.l.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.size -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.l.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, ok
	}
	return
}
