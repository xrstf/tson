package tson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Packer struct {
	types map[string]reflect.Type
}

func NewPacker() *Packer {
	return &Packer{
		types: make(map[string]reflect.Type),
	}
}

func (p *Packer) RegisterType(name string, x interface{}) {
	p.types[name] = p.ActualType(x)
}

// ActualType takes a value or a pointer to a value and returns the
// value's type.
func (p *Packer) ActualType(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

// StringifyType returns the canonical, fully-qualified type name,
// i.e. "pkg/path/to/somewhere.Type".
func (p *Packer) StringifyType(t reflect.Type) string {
	pkgPath := t.PkgPath()
	parts := strings.Split(t.Name(), ".")
	name := parts[len(parts)-1]

	return fmt.Sprintf("%s.%s", pkgPath, name)
}

type packed struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (p *Packer) Encode(value interface{}) ([]byte, error) {
	valueType := p.ActualType(value)
	typeName := p.StringifyType(valueType)
	alias := ""

	for name, t := range p.types {
		if typeName == p.StringifyType(t) {
			alias = name
		}
	}

	if alias == "" {
		return nil, fmt.Errorf("type %s has not been registered before", typeName)
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	packed := packed{
		Type:  alias,
		Value: encoded,
	}

	return json.Marshal(packed)
}

func (p *Packer) Decode(bytes []byte) (interface{}, error) {
	packed := packed{}

	err := json.Unmarshal(bytes, &packed)
	if err != nil {
		return nil, err
	}

	for name, t := range p.types {
		if name == packed.Type {
			value := reflect.New(t).Interface()
			err := json.Unmarshal(packed.Value, &value)

			return value, err
		}
	}

	return nil, fmt.Errorf("Unrecognized type '%s' found.", packed.Type)
}
