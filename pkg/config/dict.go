package config

import "sort"

// Prop represents a configuration property with its metadata.
type Prop struct {
	Name        string
	Description string
	Secret      bool
	Unique      bool
	Type        string
	DefaultVal  string
	CurrentVal  string
}

// Dict is a map of property names to Property objects.
type Dict map[string]*Prop

// Merge adds or updates a property in the dictionary.
func (d Dict) Merge(property *Prop) {
	d[property.Name] = property
}

// GetOrRegister retrieves an existing property or creates a new one if it doesn't exist.
func (d Dict) GetOrRegister(name string) *Prop {
	if p, ok := d[name]; ok {
		return p
	}

	d[name] = &Prop{Name: name}
	return d[name]
}

// Sorted returns properties sorted by name.
func (d Dict) Sorted() []*Prop {
	res := make([]*Prop, 0, len(d))
	keys := make([]string, 0, len(d))

	for k := range d {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		res = append(res, d[k])
	}

	return res
}
