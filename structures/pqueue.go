package structures

/*
This pqueue is not thread-safe.
In my case, pq is always used together with map, so lock should be put outside.
Use like this :

func Add(item) {
	mutex.Lock()
	map[item.id] = item
	pq.Push(item)
	mutex.Unlock()
}

*/

//PQItem stands for element stored in the pqueue
type PQItem struct {
	Value    interface{}
	Priority float64
	Index    int
}

//pq will shrink to save mem when meeting with two conditions when poping, see Pop below.
// 1. len(pq) < capacity(pq)/2
// 2. capacity(pq) > shrinksize
const shrinkSize = 32

//PriorityQueue implements pqueue with an array. The top item has the smallest priority.
type PriorityQueue []*PQItem

//NewPQ returns a new pqueue
func NewPQ(capacity int) PriorityQueue {
	return make(PriorityQueue, 0, capacity)
}

func (pq PriorityQueue) swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

//Push pushes item into pq
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c {
		npq := make(PriorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	item := x.(*PQItem)
	item.Index = n
	(*pq)[n] = item
	pq.up(n)
}

//Pop pops the top item from pq
func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)
	pq.swap(0, n-1)
	pq.down(0, n-1)
	if n < (c/2) && c > shrinkSize {
		npq := make(PriorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

//Peek peeks the top item
func (pq *PriorityQueue) Peek() interface{} {
	if len(*pq) == 0 {
		return nil
	}
	return (*pq)[0]
}

//Remove remove the item on position i
func (pq *PriorityQueue) Remove(i int) interface{} {
	n := len(*pq)
	if n-1 != i {
		pq.swap(i, n-1)
		pq.down(i, n-1)
		pq.up(i)
	}
	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

func (pq *PriorityQueue) up(j int) {
	for {
		i := (j - 1) / 2
		if i == j || (*pq)[j].Priority >= (*pq)[i].Priority {
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
		if j2 := j1 + 1; j2 < n && (*pq)[j1].Priority >= (*pq)[j2].Priority {
			j = j2 // = 2*i + 2  // right child
		}
		if (*pq)[j].Priority >= (*pq)[i].Priority {
			break
		}
		pq.swap(i, j)
		i = j
	}
}
