package domain

type Codec interface {
	Encode(eventName string, data any) ([]byte, error)
	Decode(payload []byte) (map[string]any, error)
}
