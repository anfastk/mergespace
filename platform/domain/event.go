package domain

type Envelope struct {
	ID      string
	Name    string
	Payload map[string]any
}
