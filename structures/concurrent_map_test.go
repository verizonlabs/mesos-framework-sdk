package structures

import (
	"reflect"
	"sync"
	"testing"
)

const mapSize = 10000

// Ensure we return the correct type and have initialized our internal data.
func TestNewConcurrentMap(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap()

	if reflect.TypeOf(m) != reflect.TypeOf(new(ConcurrentMap)) {
		t.Fatal("Creating a new concurrent map without a size gives the wrong type")
	}

	m = NewConcurrentMap(mapSize)

	if reflect.TypeOf(m) != reflect.TypeOf(new(ConcurrentMap)) {
		t.Fatal("Creating a new concurrent map with a size gives the wrong type")
	}
}

// Measures performance of creating a new ConcurrentMap.
func BenchmarkNewConcurrentMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewConcurrentMap()
	}
}

// Tests whether we can set data properly or not.
func TestConcurrentMap_Set(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap()
	var wg sync.WaitGroup

	if reflect.TypeOf(m.Set(1, 1)) != reflect.TypeOf(new(ConcurrentMap)) {
		t.Fatal("Wrong type returned from setting data")
	}

	// Goroutines are so light weight that we can easily to this.
	// Note that this barely takes up any CPU on a regular consumer machine.
	threads := 50
	wg.Add(threads)

	for i := 0; i < threads; i++ {

		// It's VERY important to close over i here.
		// If we don't, only the last value will be seen by all goroutines (usually).
		// Also, we will run into deadlocks and race conditions.
		go func(i int) {
			defer wg.Done()

			for k := 0; k < mapSize; k++ {
				m.Set(k*i, k*i)
				if m.Get(k*i) != k*i {
					t.Fatal("Could not set data in a thread-safe way")
				}
			}
		}(i)
	}

	wg.Wait()
	if m.Length() != 203459 {
		t.Fatal("Failed to properly set all data")
	}
}

// Measures performance of setting values in the map.
func BenchmarkConcurrentMap_Set(b *testing.B) {
	m := NewConcurrentMap()

	for n := 0; n < b.N; n++ {
		m.Set(n, n)
	}
}

// Tests whether we can get data properly or not.
func TestConcurrentMap_Get(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap().Set(1, 1)
	var wg sync.WaitGroup

	if reflect.TypeOf(m.Get(1)) != reflect.TypeOf(*new(int)) {
		t.Fatal("Wrong type returned from getting data")
	}

	// Goroutines are so light weight that we can easily to this.
	// Note that this barely takes up any CPU on a regular consumer machine.
	threads := 50
	wg.Add(threads)

	for i := 0; i < threads; i++ {

		// It's VERY important to close over i here.
		// If we don't, only the last value will be seen by all goroutines (usually).
		// Also, we will run into deadlocks and race conditions.
		go func(i int) {
			defer wg.Done()

			for k := mapSize; k >= 0; k-- {
				m.Set(k*i, k*i)
				if m.Get(k*i) != k*i {
					t.Fatal("Could not get data in a thread-safe way")
				}
			}
		}(i)
	}

	wg.Wait()
}

// Measures performance of getting a value from the map.
func BenchmarkConcurrentMap_Get(b *testing.B) {
	m := NewConcurrentMap()
	m.Set(0, 0)

	for n := 0; n < b.N; n++ {
		m.Get(0)
	}
}

// Tests whether we can get data properly or not.
func TestConcurrentMap_Delete(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap()
	var wg sync.WaitGroup

	// Goroutines are so light weight that we can easily to this.
	// Note that this barely takes up any CPU on a regular consumer machine.
	threads := 50
	wg.Add(threads)

	for i := 0; i < threads; i++ {

		// It's VERY important to close over i here.
		// If we don't, only the last value will be seen by all goroutines (usually).
		// Also, we will run into deadlocks and race conditions.
		go func(i int) {
			defer wg.Done()

			for k := 0; k < mapSize; k++ {
				m.Set(k*i, k*i).Delete(k * i)
			}
		}(i)
	}

	wg.Wait()
}

// Measures performance of deleting values from the map.
func BenchmarkConcurrentMap_Delete(b *testing.B) {
	m := NewConcurrentMap()

	for n := 0; n < b.N; n++ {
		m.Delete(n)
	}
}

// Tests getting the correct length of the map.
func TestConcurrentMap_Length(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap().Set(1, 1).Set(2, 2)

	if m.Length() != 2 {
		t.Fatal("Failed to get the correct length")
	}
}

// Measures performance of finding the length of the map.
func BenchmarkConcurrentMap_Length(b *testing.B) {
	m := NewConcurrentMap()

	for n := 0; n < b.N; n++ {
		m.Length()
	}
}

// Ensure that we can correctly iterate over our map.
func TestConcurrentMap_Iterate(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap()
	var wg sync.WaitGroup

	// Goroutines are so light weight that we can easily to this.
	// Note that this barely takes up any CPU on a regular consumer machine.
	threads := 50
	wg.Add(threads * 2)

	for i := 0; i < mapSize; i++ {
		m.Set(i, i)
	}

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()

			for item := range m.Iterate() {

				// Maps have no ordering.
				if item.Key.(int) >= mapSize || item.Value.(int) >= mapSize {
					t.Fatal("Failed to iterate over the map concurrently")
				}
			}
		}()
	}

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()

			for item := range m.Iterate() {

				// Make sure we are able to break out safely.
				if item.Key.(int) == 1000 {
					break
				}
			}
		}()
	}

	wg.Wait()

	if reflect.TypeOf(m.Iterate()) != reflect.TypeOf(make(<-chan Item)) {
		t.Fatal("Iterator is of the wrong type")
	}
}

// Measures performance of iterating through a map.
func BenchmarkConcurrentMap_Iterate(b *testing.B) {
	m := NewConcurrentMap()
	for i := 0; i < mapSize; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		m.Iterate()
	}
}

// Measures performance of many threads setting and reading values at the same time.
// Gives a good indicator of performance with lots of contention using a mix of read and write locks.
func BenchmarkConcurrentMap_SetRead(b *testing.B) {
	m := NewConcurrentMap()
	var wg sync.WaitGroup

	// Prime our map first.
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}

	wg.Add(100)
	b.ResetTimer()

	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()

			for i := b.N; i > 0; i-- {
				m.Set(i, i)
			}
		}()

		go func() {
			defer wg.Done()

			for i := b.N; i > 0; i-- {
				m.Get(i)
			}
		}()
	}

	wg.Wait()
}

// Measures performance of many threads setting and deleting values at the same time.
// Gives a good indicator of performance with lots of contention using write locks.
func BenchmarkConcurrentMap_SetDelete(b *testing.B) {
	m := NewConcurrentMap()
	var wg sync.WaitGroup

	// Prime our map first.
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}

	wg.Add(100)
	b.ResetTimer()

	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()

			for i := b.N; i > 0; i-- {
				m.Set(i, i)
			}
		}()

		go func() {
			defer wg.Done()

			for i := b.N; i > 0; i-- {
				m.Delete(i)
			}
		}()
	}

	wg.Wait()
}
