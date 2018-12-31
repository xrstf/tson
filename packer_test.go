package tson

import (
	"testing"
)

type structA struct {
	ID int
}

type structB struct {
	ID string
}

func TestPackingUnknownTypes(t *testing.T) {
	p := NewPacker()

	encoded, err := p.Encode(structA{})
	if err == nil {
		t.Errorf("Expected encoding to fail, but returned '%s'.", encoded)
	}
}

func TestPackingSingleStruct(t *testing.T) {
	p := NewPacker()
	p.RegisterType("a", structA{})

	data := structA{
		ID: 42,
	}

	encoded, err := p.Encode(data)
	if err != nil {
		t.Errorf("Expected encoding to succeed, but error'ed out: %v", err)
	}

	decoded, err := p.Decode(encoded)
	if err != nil {
		t.Errorf("Expected decoding to succeed, but error'ed out: %v", err)
	}

	decodedData, ok := decoded.(*structA)
	if !ok {
		t.Errorf("Expected to get a *structA, but got something else: %T", decoded)
	}

	if decodedData.ID != data.ID {
		t.Error("Decoded data does not contain original values.")
	}
}

func TestPackingMultipleTypes(t *testing.T) {
	p := NewPacker()
	p.RegisterType("a", structA{})
	p.RegisterType("b", structB{})

	dataA := structA{
		ID: 42,
	}

	dataB := structB{
		ID: "foo",
	}

	encodedA, err := p.Encode(dataA)
	if err != nil {
		t.Errorf("Expected encoding A to succeed, but error'ed out: %v", err)
	}

	encodedB, err := p.Encode(dataB)
	if err != nil {
		t.Errorf("Expected encoding B to succeed, but error'ed out: %v", err)
	}

	decodedA, err := p.Decode(encodedA)
	if err != nil {
		t.Errorf("Expected decoding A to succeed, but error'ed out: %v", err)
	}

	decodedB, err := p.Decode(encodedB)
	if err != nil {
		t.Errorf("Expected decoding B to succeed, but error'ed out: %v", err)
	}

	decodedDataA, ok := decodedA.(*structA)
	if !ok {
		t.Errorf("Expected to get a *structA, but got something else: %T", decodedA)
	}

	if decodedDataA.ID != dataA.ID {
		t.Error("Decoded data does not contain original values.")
	}

	decodedDataB, ok := decodedB.(*structB)
	if !ok {
		t.Errorf("Expected to get a *structB, but got something else: %T", decodedB)
	}

	if decodedDataB.ID != dataB.ID {
		t.Error("Decoded data does not contain original values.")
	}
}

func TestUnpackingUnknownType(t *testing.T) {
	p := NewPacker()
	p.RegisterType("a", structA{})

	data := structA{
		ID: 42,
	}

	encoded, err := p.Encode(data)
	if err != nil {
		t.Errorf("Expected encoding to succeed, but error'ed out: %v", err)
	}

	p = NewPacker()

	decoded, err := p.Decode(encoded)
	if err == nil {
		t.Errorf("Expected decoding to fail, but got %#v", decoded)
	}
}
