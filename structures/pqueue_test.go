package pqueue

import (
	"fmt"
	"math/rand"
	"testing"

	"mesos-framework-sdk/test"
)

func TestPushAndPop(t *testing.T) {
	c := 100
	pq := New(c)

	for i := 0; i < c+1; i++ {
		pq.Push(&Item{
			Value:    i,
			Priority: int64(i),
		})
	}
	test.Equal(t, c+1, len(pq))
	test.Equal(t, c*2, cap(pq))

	for i := 0; i < c+1; i++ {
		item := pq.Pop().(*Item)
		test.Equal(t, int64(i), item.Priority)
	}
	test.Equal(t, c/4, cap(pq))
}

func TestRemove(t *testing.T) {
	c := 100
	pq := New(c)

	items := make(map[string]*Item)
	for i := 0; i < c; i++ {
		p := int64(rand.Intn(100000000))
		v := fmt.Sprintf("v%d", p)
		item := &Item{
			Priority: p,
			Value:    v,
		}
		items[v] = item
		pq.Push(item)
	}

	for i := 0; i < 10; i++ {
		idx := rand.Intn((c - 1) - i)
		var f *Item
		for _, item := range items {
			if item.Index == idx {
				f = item
				break
			}
		}
		rm := pq.Remove(idx)
		test.Equal(t, fmt.Sprintf("%s", f.Value.(string)), fmt.Sprintf("%s", rm.(*Item).Value.(string)))
	}

	lastPriority := pq.Pop().(*Item).Priority
	for i := 0; i < (c - 10 - 1); i++ {
		item := pq.Pop().(*Item)
		test.Equal(t, true, lastPriority <= item.Priority)
		lastPriority = item.Priority
	}
}
