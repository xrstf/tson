# Type-safe JSON

This package implements a method to marshal and unmarshal Go structs as JSON without losing
type information. It's only useful if you intend to transmit/log somewhat arbitrary structs
and then want to unmarshal them without knowing which type you received.

## Example

Suppose you have a struct that has a fancy custom type for one of its field and this fancy
type has its own `MarshalJSON` function:

```go
type FancyType int

func (f FancyType) MarshalJSON() ([]byte, error) {
	return json.Marshal("super fancy")
}

type MyStruct struct {
	Fancyness FancyType
}

m := MyStruct{
	Fancyness: FancyType(42),
}
```

You now want to store this struct as JSON, which would be simple enough using just the
default `encoding/json` package:

```go
encoded, err := json.Marshal(m)
```

And now you want to read this data back. Again, easy enough when you know the type beforehand:

```go
result := MyStruct{}

err := json.Unmarshal(encoded, &result)
```

This package helps in cases where you **don't** know the type, but still want to make sure that
your `FancyType` is still properly unmarshalled.

## Usage

The encoder in this package works by being fed a list of types that it should know. Each type
gets a name so that the data we create is indepdent of the Go package names (we *could* in
theory just use `github.com/your/package.StructName`, but this would tie the encoded data to
the internal structure of your app and make refatorings very, very painful).

```go
import "github.com/xrstf/tson"

p := tson.NewPacker()
p.RegisterType("mystruct", MyStruct{})
```

You can now marshal and unmarshal variables of type `MyStruct` like so:

```go
encoded, err := p.Encode(MyStruct{
	Fancyness: FancyType(42),
})
// encoded = '{"type":"mystruct","value":{"Fancyness":"super fancy"}}'

decoded, err := p.Decode(encoded)
// decoded is &MyStruct{....}

asserted, ok := decoded.(*MyStruct)
if !ok {
	panic("tson package has a bug")
}

fmt.Println(asserted.Fancyness) // prints "42"
```

## Why?

Because I needed something like this. The example above is obviously bogus,
because at the end we *do* assert with a concrete type, so we could have just
as easily used the type for decoding as well. But in my usecase, I will not
assert the type later on, but continue to use the `interface{}` I get out of
`Decode()`. The important part for me is that the `FancyType` is not mangled
and read back as a string.

## License

MIT
