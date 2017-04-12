package structures

import (
	"reflect"
	"sync"
	"testing"
)

// Ensure we return the correct type and have initialized our internal data.
func TestNewConcurrentMap(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap()

	if reflect.TypeOf(m) != reflect.TypeOf(new(ConcurrentMap)) {
		t.Fatal("Creating a new concurrent map without a size gives the wrong type")
	}

	m = NewConcurrentMap(100)

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

			for k := 0; k < 10000; k++ {
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

			for k := 10000; k >= 0; k-- {
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

			for k := 0; k < 10000; k++ {
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

	for i := 0; i < 10000; i++ {
		m.Set(i, i)
	}

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()

			for item := range m.Iterate() {

				// Maps have no ordering.
				if item.Key.(int) >= 10000 || item.Value.(int) >= 10000 {
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
