package config

import (
	"fmt"
	"strings"
)

type APIKey struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (a *APIKey) FromString(raw string) (*APIKey, error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("key was not in a valid format \"name:value\"")
	}

	a.Name = parts[0]
	a.Value = parts[1]
	return a, nil
}

func (a *APIKey) String() string {
	return fmt.Sprintf("%s:%s", a.Name, a.Value)
}

type APIKeys []APIKey

func (a *APIKeys) FromStrings(input []string) (*APIKeys, error) {
	for index, raw := range input {
		key, err := new(APIKey).FromString(raw)
		if err != nil {
			return nil, fmt.Errorf("key %d was invalid: %s", index, err.Error())
		}

		*a = append(*a, *key)
	}

	return a, nil
}

func (a *APIKeys) MustFromStrings(input []string) *APIKeys {
	result, err := a.FromStrings(input)
	if err != nil {
		panic(err)
	}

	return result
}

func (a *APIKeys) String() string {
	result := make([]string, len(*a))
	for index, key := range *a {
		result[index] = key.String()
	}

	return strings.Join(result, ",")
}

func (a *APIKeys) Names() []string {
	names := make([]string, len(*a))
	for index, key := range *a {
		names[index] = key.Name
	}

	return names
}

func (a *APIKeys) KeyFor(name string) *string {
	for _, key := range *a {
		if key.Name == name {
			return &key.Value
		}
	}

	return nil
}
