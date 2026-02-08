package slice

// Map applies a function to each element and returns a new slice
func Map[T, U any](s []T, fn func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Filter filters elements based on a predicate
func Filter[T any](s []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce reduces a slice to a single value
func Reduce[T, U any](s []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range s {
		result = fn(result, v)
	}
	return result
}

// ForEach applies a function to each element
func ForEach[T any](s []T, fn func(T)) {
	for _, v := range s {
		fn(v)
	}
}

// Find finds the first element that matches a predicate
func Find[T any](s []T, fn func(T) bool) (T, bool) {
	for _, v := range s {
		if fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex finds the index of the first element that matches a predicate
func FindIndex[T any](s []T, fn func(T) bool) int {
	for i, v := range s {
		if fn(v) {
			return i
		}
	}
	return -1
}

// Contains checks if an element exists
func Contains[T comparable](s []T, item T) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

// IndexOf returns the index of an element
func IndexOf[T comparable](s []T, item T) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the last index of an element
func LastIndexOf[T comparable](s []T, item T) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == item {
			return i
		}
	}
	return -1
}

// Some checks if any element matches a predicate
func Some[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

// Every checks if all elements match a predicate
func Every[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Reverse reverses a slice
func Reverse[T any](s []T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// Unique returns unique elements
func Unique[T comparable](s []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// UniqueBy returns unique elements based on a key function
func UniqueBy[T any, K comparable](s []T, fn func(T) K) []T {
	seen := make(map[K]bool)
	result := make([]T, 0)
	for _, v := range s {
		key := fn(v)
		if !seen[key] {
			seen[key] = true
			result = append(result, v)
		}
	}
	return result
}

// Flatten flattens a 2D slice
func Flatten[T any](s [][]T) []T {
	result := make([]T, 0)
	for _, v := range s {
		result = append(result, v...)
	}
	return result
}

// FlatMap maps and flattens
func FlatMap[T, U any](s []T, fn func(T) []U) []U {
	result := make([]U, 0)
	for _, v := range s {
		result = append(result, fn(v)...)
	}
	return result
}

// Chunk splits a slice into chunks
func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	result := make([][]T, 0, (len(s)+size-1)/size)
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

// Take takes the first n elements
func Take[T any](s []T, n int) []T {
	if n <= 0 {
		return nil
	}
	if n >= len(s) {
		return s
	}
	return s[:n]
}

// TakeWhile takes elements while predicate is true
func TakeWhile[T any](s []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range s {
		if !fn(v) {
			break
		}
		result = append(result, v)
	}
	return result
}

// Drop drops the first n elements
func Drop[T any](s []T, n int) []T {
	if n <= 0 {
		return s
	}
	if n >= len(s) {
		return nil
	}
	return s[n:]
}

// DropWhile drops elements while predicate is true
func DropWhile[T any](s []T, fn func(T) bool) []T {
	for i, v := range s {
		if !fn(v) {
			return s[i:]
		}
	}
	return nil
}

// Partition partitions a slice into two based on a predicate
func Partition[T any](s []T, fn func(T) bool) ([]T, []T) {
	left := make([]T, 0)
	right := make([]T, 0)
	for _, v := range s {
		if fn(v) {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	return left, right
}

// GroupBy groups elements by a key function
func GroupBy[T any, K comparable](s []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range s {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Concat concatenates slices
func Concat[T any](slices ...[]T) []T {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	result := make([]T, 0, total)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// Clone clones a slice
func Clone[T any](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	return result
}

// First returns the first element
func First[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[0], true
}

// Last returns the last element
func Last[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[len(s)-1], true
}

// At returns the element at index (supports negative indices)
func At[T any](s []T, index int) (T, bool) {
	if index < 0 {
		index = len(s) + index
	}
	if index < 0 || index >= len(s) {
		var zero T
		return zero, false
	}
	return s[index], true
}

// IsEmpty checks if slice is empty
func IsEmpty[T any](s []T) bool {
	return len(s) == 0
}

// Len returns the length
func Len[T any](s []T) int {
	return len(s)
}
