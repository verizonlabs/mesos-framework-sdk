package structures

import (
	"sync"
)

//Item stands for element stored in the pqueue
type Item struct {
	Value    interface{}
	Priority int64
	Index    int
}

//pq will shrink to save mem if meeting two conditions when poping, see Pop below.
// 1. len(pq) < capacity(pq)/2
// 2. capacity(pq) > shrinksize
const shrinkSize = 32

//Items is a slice of Item
type Items []*Item

//PriorityQueue implements pqueue with an array. The top item has the smallest priority.
type PriorityQueue struct {
	lock sync.RWMutex
	Items
}

//New returns a new pqueue
func New(capacity int) PriorityQueue {
	return PriorityQueue{
		Items: make(Items, 0, capacity),
	}
}

// Len returns size of items in pq
func (pq *PriorityQueue) Len() int {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return len(pq.Items)
}

func (pq *PriorityQueue) swap(i, j int) {
	pq.Items[i], pq.Items[j] = pq.Items[j], pq.Items[i]
	pq.Items[i].Index = i
	pq.Items[j].Index = j
}

// Cap returns capacity of pq
func (pq *PriorityQueue) Cap() int {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return cap(pq.Items)
}

//Push pushes item into pq
func (pq *PriorityQueue) Push(x interface{}) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	n := len(pq.Items)
	c := cap(pq.Items)
	if n+1 > c {
		newItems := make(Items, n, c*2)
		copy(newItems, pq.Items)
		pq.Items = newItems
	}
	pq.Items = (pq.Items)[0 : n+1]
	item := x.(*Item)
	item.Index = n
	(pq.Items)[n] = item
	pq.up(n)
}

//Pop pops the top item from pq
func (pq *PriorityQueue) Pop() interface{} {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	n := len(pq.Items)
	c := cap(pq.Items)
	pq.swap(0, n-1)
	pq.down(0, n-1)
	if n < (c/2) && c > shrinkSize {
		newItems := make(Items, n, c/2)
		copy(newItems, pq.Items)
		pq.Items = newItems
	}
	item := (pq.Items)[n-1]
	item.Index = -1
	pq.Items = (pq.Items)[0 : n-1]
	return item
}

//Peek peeks the top item
func (pq *PriorityQueue) Peek() interface{} {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	if len(pq.Items) == 0 {
		return nil
	}
	return (pq.Items)[0]
}

//Remove remove the item on position i
func (pq *PriorityQueue) Remove(i int) interface{} {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	n := len(pq.Items)
	if n-1 != i {
		pq.swap(i, n-1)
		pq.down(i, n-1)
		pq.up(i)
	}
	item := pq.Items[n-1]
	item.Index = -1
	pq.Items = pq.Items[0 : n-1]
	return item
}

func (pq *PriorityQueue) up(j int) {
	for {
		i := (j - 1) / 2
		if i == j || pq.Items[j].Priority >= pq.Items[i].Priority {
			break
		}
		pq.swap(i, j)
		j = i
	}
}

func (pq *PriorityQueue) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && pq.Items[j1].Priority >= pq.Items[j2].Priority {
			j = j2 // = 2*i + 2  // right child
		}
		if pq.Items[j].Priority >= pq.Items[i].Priority {
			break
		}
		pq.swap(i, j)
		i = j
	}
}
