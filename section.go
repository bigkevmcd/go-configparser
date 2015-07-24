package configparser

import "strings"

type Section struct {
	Name    string
	options Dict
	lookup  Dict
}

func (s *Section) Add(key, value string) error {
	lookupKey := s.safeKey(key)

	s.options[key] = s.safeValue(value)
	s.lookup[lookupKey] = key
	return nil
}

func (s *Section) Get(key string) (string, error) {
	lookupKey, present := s.lookup[s.safeKey(key)]
	if !present {
		return "", getNoOptionError(s.Name, key)
	}
	if value, present := s.options[lookupKey]; present {
		return value, nil
	}
	return "", getNoOptionError(s.Name, key)
}

func (s *Section) Options() []string {
	return s.options.Keys()
}

func (s *Section) Items() Dict {
	return s.options
}

func (s *Section) safeValue(in string) string {
	// Same as safeKey for now.
	return s.safeKey(in)
}

func (s *Section) safeKey(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}

func (s *Section) Remove(key string) error {
	_, present := s.options[key]
	if !present {
		return getNoOptionError(s.Name, key)
	}

	// delete doesn't return anything, but this does require
	// that the passed key to be removed matches the options key.
	delete(s.lookup, s.safeKey(key))
	delete(s.options, key)
	return nil
}

func newSection(name string) *Section {
	return &Section{
		Name:    name,
		options: make(Dict),
		lookup:  make(Dict),
	}
}
