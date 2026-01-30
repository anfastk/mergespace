package avro

import "github.com/riferrei/srclient"

type Registry struct {
	client *srclient.SchemaRegistryClient
}

func NewRegistry(url string) *Registry {
	return &Registry{
		client: srclient.CreateSchemaRegistryClient(url),
	}
}

func (r *Registry) Register(subject string, schema string) (int, error) {
	s, err := r.client.CreateSchema(subject, schema, srclient.Avro)
	if err != nil {
		return 0, err
	}
	return s.ID(), nil
}

func (r *Registry) Get(id int) (string, error) {
	s, err := r.client.GetSchema(id)
	if err != nil {
		return "", err
	}
	return s.Schema(), nil
}
