package maputil

// Keys returns the keys of a map
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values of a map
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Entries returns key-value pairs
func Entries[K comparable, V any](m map[K]V) []struct {
	Key   K
	Value V
} {
	entries := make([]struct {
		Key   K
		Value V
	}, 0, len(m))
	for k, v := range m {
		entries = append(entries, struct {
			Key   K
			Value V
		}{Key: k, Value: v})
	}
	return entries
}

// FromEntries creates a map from key-value pairs
func FromEntries[K comparable, V any](entries []struct {
	Key   K
	Value V
}) map[K]V {
	m := make(map[K]V, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// Map transforms values using a function
func Map[K comparable, V, U any](m map[K]V, fn func(V) U) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = fn(v)
	}
	return result
}

// MapKeys transforms keys using a function
func MapKeys[K, K2 comparable, V any](m map[K]V, fn func(K) K2) map[K2]V {
	result := make(map[K2]V, len(m))
	for k, v := range m {
		result[fn(k)] = v
	}
	return result
}

// Filter filters entries based on a predicate
func Filter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}

// ForEach iterates over entries
func ForEach[K comparable, V any](m map[K]V, fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// Find finds an entry based on a predicate
func Find[K comparable, V any](m map[K]V, fn func(K, V) bool) (K, V, bool) {
	for k, v := range m {
		if fn(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Some checks if any entry matches a predicate
func Some[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if fn(k, v) {
			return true
		}
	}
	return false
}

// Every checks if all entries match a predicate
func Every[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

// ContainsKey checks if a key exists
func ContainsKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// ContainsValue checks if a value exists
func ContainsValue[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

// Get gets a value with a default
func Get[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

// GetOrPut gets a value or inserts a default
func GetOrPut[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	m[key] = defaultValue
	return defaultValue
}

// Merge merges multiple maps (later maps override earlier)
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	total := 0
	for _, m := range maps {
		total += len(m)
	}
	result := make(map[K]V, total)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Clone clones a map
func Clone[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Delete removes keys matching a predicate
func Delete[K comparable, V any](m map[K]V, fn func(K, V) bool) {
	for k, v := range m {
		if fn(k, v) {
			delete(m, k)
		}
	}
}

// Clear clears a map
func Clear[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

// IsEmpty checks if map is empty
func IsEmpty[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

// Len returns the length
func Len[K comparable, V any](m map[K]V) int {
	return len(m)
}

// Invert inverts a map (values become keys)
func Invert[K, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// GroupBy groups values by a key function
func GroupBy[K comparable, V any, T any](s []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range s {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// CountBy counts occurrences by a key function
func CountBy[K comparable, V any](s []V, fn func(V) K) map[K]int {
	result := make(map[K]int)
	for _, v := range s {
		key := fn(v)
		result[key]++
	}
	return result
}
