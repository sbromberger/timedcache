# timedmap
## A generic k/v store with expiration

Package `timedmap` provides a `Map` object that is similar to Go's
native `map` but supports expiration of the entries in the map.
`Map`s are fully generic and will support any values that do not contain
mutexes.

The following example code demonstrates the core concepts:

```go

// create a new timedmap with a 3 second default timeout for entries,
// ints for keys, and strings for values
myMap := timedmap.New[uint32, string](time.Duration(3 * time.Second))

// Set() adds entries to the map
a.Set(10, "five plus five") 
a.Set(20, "twenty")
a.Set(10, "ten")

// Get() retrieves a value by key
fmt.Println(a.Get(10))  // "ten"

// Dump() returns a standard Go map containing unexpired entries
fmt.Println(a.Dump())

time.Sleep(2 * time.Second)
// before entry "20: twenty" expires, reset its timer to the default
a.Reset(20)

time.Sleep(2 * time.Second)
// "10: ten" has expired by now
fmt.Println(a.Dump())

// Purge manually removes all expired entries from the map
fmt.Printf("%d entries purged\n", a.Purge())

// Manually set an expiration for entry with key 20
a.SetExpiration(20, time.Duration(100 * time.Second))

// Delete entry with key 20
a.Delete(20)
```
