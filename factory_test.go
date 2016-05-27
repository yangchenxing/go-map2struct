package map2struct

import (
	"fmt"
	"reflect"
	"testing"
)

type Stringer interface {
	String() string
}

type Foo struct {
	Text string
}

func (foo Foo) String() string {
	return "foo:" + foo.Text
}

type Bar struct {
	Int int
}

func (bar Bar) String() string {
	return fmt.Sprintf("bar:%d", bar.Int)
}

func initializeInstance(i interface{}) error {
	return nil
}

func TestFactory(t *testing.T) {
	factory := NewGeneralInterfaceFactory(reflect.TypeOf((*Stringer)(nil)).Elem(),
		"type", initializeInstance)
	factory.RegisterType("Foo", reflect.TypeOf((*Foo)(nil)).Elem())
	factory.RegisterType("Bar", reflect.TypeOf((*Bar)(nil)).Elem())
	RegisterFactory(factory)
	src := map[string]interface{}{
		"Stringer1": map[string]interface{}{
			"type": "Foo",
			"Text": "hello",
		},
		"Stringer2": map[string]interface{}{
			"type": "Bar",
			"Int":  1,
		},
	}
	type TestStruct struct {
		Stringer1 Stringer
		Stringer2 Stringer
	}
	var output TestStruct
	if err := UnmarshalMap(&output, src); err != nil {
		t.Error("unmarshal map fail:", err.Error())
		return
	}
	if output.Stringer1.String() != "foo:hello" || output.Stringer2.String() != "bar:1" {
		t.Error("unexpected output:", output)
		return
	}

	type TestStringer Stringer
	src = map[string]interface{}{
		"type": "Foo",
		"Text": "world",
	}
	var tsOutput Stringer
	if err := UnmarshalMap(&tsOutput, src); err != nil {
		t.Error("unmarshal map fail:", err.Error())
		return
	}
	if tsOutput.String() != "foo:world" {
		t.Error("unexpected output:", tsOutput)
		return
	}

	// test unknown type
	src = map[string]interface{}{
		"type": "unknown",
	}
	if err := UnmarshalMap(&tsOutput, src); err == nil {
		t.Error("unexpected unmarshal success:", tsOutput)
		return
	}
}
