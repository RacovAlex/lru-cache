package cache

import (
	"container/list"
	"sync"
)

// LRUCache - структура для реализации LRU-кэша, реализующая интерфейс Cache.
// Каждый элемент списка хранит указатель на структуру entry, которая содержит ключ и значение кэша.
type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*list.Element
	list     *list.List
	mu       sync.Mutex
}

// entry - ключ-значение для хранения в списке, требуется для удобного удаления элементов из cache.
type entry[K comparable, V any] struct {
	key   K
	value V
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element, capacity),
		list:     list.New(),
	}
}

// Get возвращает значение по ключу и флаг его наличия и перемещает элемент в начало (если он в кэше).
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		// Поскольку Value является типом interface{}, нам нужно выполнить приведение типа
		return element.Value.(*entry[K, V]).value, ok
	}
	var zeroValue V
	return zeroValue, false
}

// Put добавляет или обновляет элемент в кэше.
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Если элемент в кэше - обновляем его значение и перемещаем в начало списка
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		element.Value.(*entry[K, V]).value = value
		return
	}
	// Если длина листа превышает емкость кэша - удаляем самый старый элемент
	if c.list.Len() >= c.capacity {
		c.removeOldest()
	}
	// Помещаем элемент в начало списка и обновляем map cache
	element := c.list.PushFront(&entry[K, V]{key: key, value: value})
	c.cache[key] = element
}

// removeOldest удаляет самый старый элемент.
func (c *LRUCache[K, V]) removeOldest() {
	element := c.list.Back()
	if element != nil {
		c.list.Remove(element)
		delete(c.cache, element.Value.(*entry[K, V]).key)
	}
}
