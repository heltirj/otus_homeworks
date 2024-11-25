package storage

import "fmt"

type Type int

const (
	SQLStorage Type = 1
	InMemory   Type = 2
)

var storageTypeStrings = map[string]Type{
	"sql":       SQLStorage,
	"in-memory": InMemory,
}

func (s *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var storageType string
	if err := unmarshal(&storageType); err != nil {
		return err
	}
	*s = storageTypeStrings[storageType]
	if *s == 0 {
		return fmt.Errorf("invalid storage type: %s", storageType)
	}

	return nil
}
