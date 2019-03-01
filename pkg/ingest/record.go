package ingest

type Record interface {
	Get(key string) interface{}
	Put(key string, value interface{})
}

type MapRecord struct {
	fields map[string]interface{}
}

func (c *MapRecord) Get(key string) interface{} {
	return c.fields[key]
}

func (c *MapRecord) Put(key string, value interface{}) {
	c.fields[key] = value
}

func NewMapRecord() *MapRecord {
	return &MapRecord{fields: make(map[string]interface{})}
}
