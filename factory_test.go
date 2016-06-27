package map2struct

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type stringer interface {
	String() string
}

type texter interface {
	Text() string
}

type foo struct {
	Text string
}

func (foo foo) String() string {
	return "foo:" + foo.Text
}

type bar struct {
	Duration time.Duration
}

func (bar bar) String() string {
	return fmt.Sprintf("bar:%s", bar.Duration)
}

func initializeInstance(i interface{}) error {
	return nil
}

func initializeInstanceFail(i interface{}) error {
	return fmt.Errorf("error")
}

func TestRegisterFactory(t *testing.T) {
	defer func() { factories = make(map[string]Factory) }()
	var s string
	factory := NewGeneralInterfaceFactory(reflect.TypeOf(s), "type", nil)
	RegisterFactory(factory)
	if f := factories[getTypeName(reflect.TypeOf(s))]; f != factory {
		t.Error("register factory fail:", f)
		return
	}
}

func TestGeneralInterfaceFactory(t *testing.T) {
	defer func() { factories = make(map[string]Factory) }()
	factory := NewGeneralInterfaceFactory(reflect.TypeOf((*stringer)(nil)).Elem(),
		"type", initializeInstance)
	factory.RegisterType("Foo", reflect.TypeOf((*foo)(nil)).Elem())
	factory.RegisterType("Bar", reflect.TypeOf((*bar)(nil)).Elem())
	factory.RegisterInstance("Instance", &foo{Text: "instance"})
	RegisterFactory(factory)
	src := map[string]interface{}{
		"Stringer1": map[string]interface{}{
			"type": "Foo",
			"Text": "hello",
		},
		"Stringer2": map[string]interface{}{
			"type":     "Bar",
			"Duration": "1s",
		},
		"Instance": map[string]interface{}{
			"type": "Instance",
		},
	}
	type TestStruct struct {
		Stringer1 stringer
		Stringer2 stringer
		Instance  stringer
	}
	type TestStruct2 struct {
		Texter texter
	}
	var output TestStruct
	if err := Unmarshal(&output, src); err != nil {
		t.Error("unmarshal map fail:", err.Error())
		return
	}
	if output.Stringer1.String() != "foo:hello" || output.Stringer2.String() != "bar:1s" || output.Instance.String() != "foo:instance" {
		t.Error("unexpected output:", output)
		return
	}

	src = map[string]interface{}{
		"type": "Foo",
		"Text": "world",
	}
	var s stringer
	var x texter
	if err := Unmarshal(&s, src); err != nil {
		t.Error("unmarshal map fail:", err.Error())
		return
	}
	if s.String() != "foo:world" {
		t.Error("unexpected output:", s)
		return
	}

	// test unknown type
	src = map[string]interface{}{
		"type": "unknown",
	}
	if err := Unmarshal(&s, src); err == nil {
		t.Error("unexpected unmarshal success:", s)
		return
	}

	// test missing type key
	src = map[string]interface{}{
		"Text": "hello",
	}
	if err := Unmarshal(&s, src); err == nil {
		t.Error("unexpected unmarshal success:", s)
		return
	}
	// test unknown dest type
	if err := Unmarshal(&x, src); err == nil {
		t.Error("unexpected unmarshal success:", t)
		return
	}

	// test ptr
	src = map[string]interface{}{
		"type": "Foo",
		"Text": "world",
	}
	factory = NewGeneralInterfaceFactory(reflect.TypeOf((*stringer)(nil)).Elem(),
		"type", initializeInstance)
	factory.RegisterType("Foo", reflect.TypeOf((**foo)(nil)).Elem())
	factory.RegisterType("Bar", reflect.TypeOf((**bar)(nil)).Elem())
	RegisterFactory(factory)
	if err := Unmarshal(&s, src); err != nil {
		t.Error("unmarshal map fail:", err.Error())
		return
	}
	if s.String() != "foo:world" {
		t.Error("unexpected output:", s)
		return
	}

	// test initializer fail
	factory = NewGeneralInterfaceFactory(reflect.TypeOf((*stringer)(nil)).Elem(),
		"type", initializeInstanceFail)
	factory.RegisterType("Foo", reflect.TypeOf((**foo)(nil)).Elem())
	factory.RegisterType("Bar", reflect.TypeOf((**bar)(nil)).Elem())
	RegisterFactory(factory)
	if err := Unmarshal(&s, src); err == nil {
		t.Error("unexpected unmarshal success:", s)
		return
	}

	// test interval fail
	src = map[string]interface{}{
		"type":     "Bar",
		"Duration": "s",
	}
	if err := Unmarshal(&s, src); err == nil {
		t.Error("unexpected unmarshal success:", s)
		return
	}
}
