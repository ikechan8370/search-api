package serp

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

// no global lock

// CacheItem represents a single cache item with an expiration time and value.
type CacheItem struct {
	Value      interface{} // The cached value
	ExpireTime time.Time   // The time when the cache item will expire
}

// Cache is the structure that holds the cache data and associated methods.
type Cache struct {
	data map[string][]CacheItem // Stores cached data where key is the cache key and value is a list of cache items
	mu   sync.RWMutex           // Mutex to handle concurrent access to the cache
	n    int                    // The maximum number of cache copies per key
}

// NewCache creates a new Cache instance with the given capacity for cache copies.
func NewCache(n int) *Cache {
	return &Cache{
		data: make(map[string][]CacheItem),
		n:    n,
	}
}

// Set adds a new cache item for a specific key, if there is space for more cache copies.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get the current list of cache items for the given key
	cacheItems, exists := c.data[key]
	if !exists {
		cacheItems = []CacheItem{} // If no cache exists for the key, initialize an empty list
	}

	// If the cache list is already full, discard the oldest item (this could be adjusted to any eviction strategy)
	if len(cacheItems) >= c.n {
		cacheItems = cacheItems[1:]
	}

	// Add the new cache item with a 6-hour expiration time
	cacheItems = append(cacheItems, CacheItem{
		Value:      value,
		ExpireTime: time.Now().Add(6 * time.Hour), // 6-hour expiration
	})

	// Update the cache with the new list of items for the key
	c.data[key] = cacheItems
}

// Get retrieves a random valid cache item for a specific key.
// It removes expired items to avoid memory leaks.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()         // Lock for write access to modify the cache
	defer c.mu.Unlock() // Ensure the lock is released when the method finishes

	// Retrieve the list of cached items for the given key
	cacheItems, exists := c.data[key]
	if !exists || len(cacheItems) < c.n {
		return nil, false // If no cache for the key or not enough items, return a cache miss
	}

	// Filter out expired cache items and remove them from the list
	validCacheItems := []CacheItem{}
	for _, item := range cacheItems {
		// Check if the item has not expired
		if time.Now().Before(item.ExpireTime) {
			validCacheItems = append(validCacheItems, item) // Add to valid items if not expired
		}
	}

	// If all cache items have expired, return a cache miss
	if len(validCacheItems) == 0 {
		delete(c.data, key) // Remove the key from the cache if all items are expired
		return nil, false
	}

	// Update the cache with only valid items (remove expired ones)
	c.data[key] = validCacheItems

	// Randomly pick one valid cache item from the list
	randomIndex := rand.Intn(len(validCacheItems))
	return validCacheItems[randomIndex].Value, true
}

// Delete method to remove cache based on key or pattern
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Special case for "all" - remove everything
	if key == "all" {
		c.data = make(map[string][]CacheItem) // Reset the cache
		return
	}

	// Loop through all keys and delete those that match the given pattern
	for k := range c.data {
		if strings.HasPrefix(k, key) {
			delete(c.data, k) // Delete matching key
		}
	}
}
