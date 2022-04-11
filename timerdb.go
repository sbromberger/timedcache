// Package timerdb provides a key/value store called a `timerdb.Map` that supports "aging out"
// or expiration of its entries based on a timeout value set at `Map` creation.
// Values are set and retrieved via `Set()` and `Get` methods, respectively. If an entry has
// timed out, it is not retrievable.
// Known restrictions include the following:
//
// - there is no way to iterate (range) over the contents of the `Map`.
//
// - using values that contain mutexes is not safe, as values may internally be copied.
package timerdb

import (
	"sync"
	"time"
)

// mapEntry is a wrapper around a generic value V, adding
// an expiration field.
type mapEntry[V any] struct {
	expiration time.Time
	v          V
}

// isExpired returns true if the expiration time of the mapEntry
// has passed, otherwise false.
func (me *mapEntry[V]) isExpired() bool {
	return time.Now().After(me.expiration)
}

// Map is a generic key/value store that expires entries after a
// user-defined timeout period.
type Map[K comparable, V any] struct {
	timeout time.Duration
	mu      sync.RWMutex
	m       map[K]mapEntry[V]
}

// New creates a new Map.
func New[K comparable, V any](t time.Duration) *Map[K, V] {
	m := map[K]mapEntry[V]{}
	return &Map[K, V]{m: m, timeout: t}
}

// Get gets a value by key from a Map along with a boolean indicating
// whether the key was found. If not found, the value will be the zero
// value.
func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	me, found := m.m[k]
	if !found || me.isExpired() {
		var zero V
		return zero, false
	}
	return me.v, true
}

// Set sets a key/value pair in a map along with a timeout if it does not already exist.
// If the entry exists, Set() will reset the timer value.
func (m *Map[K, V]) Set(k K, v V) {
	expires := time.Now().Add(m.timeout)
	m.mu.Lock()
	m.m[k] = mapEntry[V]{v: v, expiration: expires}
	m.mu.Unlock()
}

// Delete deletes an entry from a Map given its key. If the key does not
// exist in the Map, the function returns false and does nothing. If the entry
// exists but is expired, the function will return false and the entry will
// be removed. If the delete is successful, the function will return true.
func (m *Map[K, V]) Delete(k K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, found := m.m[k]
	retval := found && !entry.isExpired()
	delete(m.m, k) // this is a nop if !found, so safe to do here.
	return retval
}

// SetExpiration sets a custom expiration time for a Map entry given its key. If
// the key does not exist in the Map, the function returns false and does nothing.
func (m *Map[K, V]) SetExpiration(k K, expires time.Time) bool {
	m.mu.RLock()
	v, found := m.m[k]
	m.mu.RUnlock()
	if !found {
		return false
	}
	v.expiration = expires
	m.mu.Lock()
	m.m[k] = v
	m.mu.Unlock()
	return true
}

// Reset resets the timeout for a Map entry given its key. If the key
// does not exist in the Map, the function returns false and does nothing.
func (m *Map[K, V]) Reset(k K) bool {
	return m.SetExpiration(k, time.Now().Add(m.timeout))
}

// Purge deletes all expired entries from the Map and returns
// the number of deleted entries.
// Purge should be called sparingly as it locks the Map
// while it iterates over it.
func (m *Map[K, V]) Purge() int {
	var i int
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.m {
		if v.isExpired() {
			i++
			delete(m.m, k)
		}
	}
	return i
}
