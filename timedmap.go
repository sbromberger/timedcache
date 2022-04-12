// Package timedmap provides a key/value store called a `timedmap.Map` that supports "aging out"
// of its entries based on a defaultExpiration value set at `Map` creation.
// Values are set and retrieved via `Set()` and `Get()` methods, respectively. If an entry has
// timed out, it is not retrievable. Known restrictions include the following:
//
// - there is no way to iterate (range) over the contents of the `Map`.
//
// - using values that contain mutexes is not safe, as values may be copied in the internal methods.
package timedmap

import (
	"sync"
	"time"
)

// mapEntry is a wrapper around a generic value V, adding
// an expiresAt field.
type mapEntry[V any] struct {
	expiresAt time.Time
	v         V
}

// isExpired returns true if the expiration time of the mapEntry
// has passed, otherwise false.
func (me *mapEntry[V]) isExpired() bool {
	return time.Now().After(me.expiresAt)
}

// Map is a generic key/value store that expires entries after a
// user-defined defaultExpiration period.
type Map[K comparable, V any] struct {
	defaultExpiration time.Duration
	mu                sync.RWMutex
	m                 map[K]mapEntry[V]
}

// New creates a new Map.
func New[K comparable, V any](t time.Duration) *Map[K, V] {
	m := map[K]mapEntry[V]{}
	return &Map[K, V]{m: m, defaultExpiration: t}
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

// Set sets a key/value pair in a map (along with a defaultExpiration) if it
// does not already exist. If the entry exists, Set() will reset the timer value.
func (m *Map[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[k] = mapEntry[V]{v: v, expiresAt: time.Now().Add(m.defaultExpiration)}
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
	m.mu.Lock()
	defer m.mu.Unlock()
	v, found := m.m[k]
	if !found {
		return false
	}
	v.expiresAt = expires
	m.m[k] = v
	return true
}

// Reset resets the expiration for a Map entry given its key. If the key
// does not exist in the Map, the function returns false and does nothing.
func (m *Map[K, V]) Reset(k K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, found := m.m[k]
	if !found {
		return false
	}
	v.expiresAt = time.Now().Add(m.defaultExpiration)
	m.m[k] = v
	return true
}

// Purge deletes all expired entries from the Map and returns
// the number of deleted entries.
// Purge should be called sparingly as it locks the Map
// for the duration of the iteration.
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

// Dump returns a standard map containing the unexpired
// values within the Map.
func (m *Map[K, V]) Dump() map[K]V {
	dumpm := map[K]V{}
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.m {
		if !v.isExpired() {
			dumpm[k] = v.v
		}
	}
	return dumpm
}
