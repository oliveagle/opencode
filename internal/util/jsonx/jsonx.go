package jsonx

import (
	"encoding/json"
	"io"
	"os"
)

// Marshal marshals to JSON
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent marshals to indented JSON
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal unmarshals from JSON
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MarshalString marshals to JSON string
func MarshalString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalString unmarshals from JSON string
func UnmarshalString(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

// Encode encodes to writer
func Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// Decode decodes from reader
func Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// ReadFile reads JSON from file
func ReadFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteFile writes JSON to file
func WriteFile(path string, v interface{}, perm os.FileMode) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

// MustMarshal marshals or panics
func MustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// MustMarshalString marshals to string or panics
func MustMarshalString(v interface{}) string {
	s, err := MarshalString(v)
	if err != nil {
		panic(err)
	}
	return s
}

// MustUnmarshal unmarshals or panics
func MustUnmarshal(data []byte, v interface{}) {
	if err := json.Unmarshal(data, v); err != nil {
		panic(err)
	}
}

// Map is a type alias for JSON object
type Map = map[string]interface{}

// Array is a type alias for JSON array
type Array = []interface{}

// ToMap converts to map
func ToMap(v interface{}) (Map, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m Map
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// FromMap converts from map
func FromMap(m Map, v interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// Merge merges multiple JSON maps
func Merge(maps ...Map) Map {
	result := make(Map)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Clone clones a JSON map
func Clone(m Map) Map {
	result := make(Map)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Get gets a value from nested map
func Get(m Map, keys ...string) (interface{}, bool) {
	for i, key := range keys {
		if i == len(keys)-1 {
			v, ok := m[key]
			return v, ok
		}
		next, ok := m[key].(Map)
		if !ok {
			return nil, false
		}
		m = next
	}
	return nil, false
}

// GetString gets a string value
func GetString(m Map, keys ...string) (string, bool) {
	v, ok := Get(m, keys...)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetInt gets an int value
func GetInt(m Map, keys ...string) (int, bool) {
	v, ok := Get(m, keys...)
	if !ok {
		return 0, false
	}
	// JSON numbers are float64
	if f, ok := v.(float64); ok {
		return int(f), true
	}
	i, ok := v.(int)
	return i, ok
}

// GetFloat gets a float64 value
func GetFloat(m Map, keys ...string) (float64, bool) {
	v, ok := Get(m, keys...)
	if !ok {
		return 0, false
	}
	f, ok := v.(float64)
	return f, ok
}

// GetBool gets a bool value
func GetBool(m Map, keys ...string) (bool, bool) {
	v, ok := Get(m, keys...)
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// GetArray gets an array value
func GetArray(m Map, keys ...string) (Array, bool) {
	v, ok := Get(m, keys...)
	if !ok {
		return nil, false
	}
	a, ok := v.(Array)
	return a, ok
}

// GetMap gets a nested map value
func GetMap(m Map, keys ...string) (Map, bool) {
	v, ok := Get(m, keys...)
	if !ok {
		return nil, false
	}
	nm, ok := v.(Map)
	return nm, ok
}

// Set sets a value in nested map
func Set(m Map, value interface{}, keys ...string) {
	for i, key := range keys {
		if i == len(keys)-1 {
			m[key] = value
			return
		}
		next, ok := m[key].(Map)
		if !ok {
			next = make(Map)
			m[key] = next
		}
		m = next
	}
}

// Delete deletes a value from nested map
func Delete(m Map, keys ...string) {
	for i, key := range keys {
		if i == len(keys)-1 {
			delete(m, key)
			return
		}
		next, ok := m[key].(Map)
		if !ok {
			return
		}
		m = next
	}
}

// Pretty returns pretty-printed JSON
func Pretty(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Compact returns compact JSON
func Compact(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
