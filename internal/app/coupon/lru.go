package coupon

import "sync"

type lruNode struct {
	key   string
	value bool
	prev  *lruNode
	next  *lruNode
}

type LRUCache struct {
	cap   int
	mu    sync.Mutex
	items map[string]*lruNode
	head  *lruNode
	tail  *lruNode
}

func NewLRU(cap int) *LRUCache {
	return &LRUCache{
		cap:   cap,
		items: make(map[string]*lruNode, cap),
	}
}

func (l *LRUCache) Get(k string) (bool, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if n, ok := l.items[k]; ok {
		l.moveToFront(n)
		return n.value, true
	}
	return false, false
}

func (l *LRUCache) Set(k string, v bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if n, ok := l.items[k]; ok {
		n.value = v
		l.moveToFront(n)
		return
	}
	n := &lruNode{key: k, value: v}
	l.items[k] = n
	l.addToFront(n)
	if len(l.items) > l.cap {
		l.removeOldest()
	}
}

func (l *LRUCache) addToFront(n *lruNode) {
	n.prev = nil
	n.next = l.head
	if l.head != nil {
		l.head.prev = n
	}
	l.head = n
	if l.tail == nil {
		l.tail = n
	}
}

func (l *LRUCache) moveToFront(n *lruNode) {
	if l.head == n {
		return
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	if l.tail == n {
		l.tail = n.prev
	}
	n.prev = nil
	n.next = l.head
	if l.head != nil {
		l.head.prev = n
	}
	l.head = n
}

func (l *LRUCache) removeOldest() {
	if l.tail == nil {
		return
	}
	old := l.tail
	if old.prev != nil {
		old.prev.next = nil
	}
	l.tail = old.prev
	delete(l.items, old.key)
}
