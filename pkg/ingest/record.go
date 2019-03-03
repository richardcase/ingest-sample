package ingest

// Record represents a generic key-value record
type Record interface {
	Get(key string) interface{}
	Put(key string, value interface{})
}

// MapRecord is a key-value record that is backed by a map
type MapRecord struct {
	fields map[string]interface{}
}

// Get is get the value for a given key
func (c *MapRecord) Get(key string) interface{} {
	return c.fields[key]
}

// Put will store a value for a specific key
func (c *MapRecord) Put(key string, value interface{}) {
	c.fields[key] = value
}

// NewMapRecord creates a new MapRecord
func NewMapRecord() *MapRecord {
	return &MapRecord{fields: make(map[string]interface{})}
}
