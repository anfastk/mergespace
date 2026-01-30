package avro

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/hamba/avro/v2"
)

const magicByte = 0

type Codec struct {
	registry *Registry

	mu      sync.RWMutex
	ids     map[string]int
	schemas map[int]avro.Schema
}

func NewCodec(registry *Registry) *Codec {
	return &Codec{
		registry: registry,
		ids:      make(map[string]int),
		schemas:  make(map[int]avro.Schema),
	}
}

func (c *Codec) Register(
	eventName string,
	subject string,
	schemaStr string,
) error {

	id, err := c.registry.Register(subject, schemaStr)
	if err != nil {
		return err
	}

	parsed, err := avro.Parse(schemaStr)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.ids[eventName] = id
	c.schemas[id] = parsed
	c.mu.Unlock()

	return nil
}

func (c *Codec) Encode(eventName string, data any) ([]byte, error) {
	c.mu.RLock()
	id, ok := c.ids[eventName]
	schema := c.schemas[id]
	c.mu.RUnlock()

	if !ok {
		return nil, errors.New("schema not registered for event")
	}

	bin, err := avro.Marshal(schema, data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteByte(magicByte)
	_ = binary.Write(&buf, binary.BigEndian, int32(id))
	buf.Write(bin)

	return buf.Bytes(), nil
}

func (c *Codec) Decode(payload []byte) (map[string]any, error) {
	if len(payload) < 5 || payload[0] != magicByte {
		return nil, errors.New("invalid avro payload")
	}

	id := int(binary.BigEndian.Uint32(payload[1:5]))

	schema, err := c.getSchema(id)
	if err != nil {
		return nil, err
	}

	var out map[string]any
	if err := avro.Unmarshal(schema, payload[5:], &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Codec) getSchema(id int) (avro.Schema, error) {
	c.mu.RLock()
	if s, ok := c.schemas[id]; ok {
		c.mu.RUnlock()
		return s, nil
	}
	c.mu.RUnlock()

	raw, err := c.registry.Get(id)
	if err != nil {
		return nil, err
	}

	parsed, err := avro.Parse(raw)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.schemas[id] = parsed
	c.mu.Unlock()

	return parsed, nil
}
