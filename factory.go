package map2struct

import (
	"fmt"
	"reflect"
)

type Factory interface {
	GetInstanceType() reflect.Type
	Create(map[string]interface{}) (interface{}, error)
}

var (
	factories = make(map[string]Factory)
)

func RegisterFactory(factory Factory) {
	if factory != nil {
		factories[getTypeName(factory.GetInstanceType())] = factory
	}
}

func CreateByFactory(typ reflect.Type, data map[string]interface{}) (interface{}, error) {
	typeName := getTypeName(typ)
	if factory := factories[typeName]; factory != nil {
		return factory.Create(data)
	}
	return nil, fmt.Errorf("unregistered type: %q", typeName)
}

func getTypeName(typ reflect.Type) string {
	if pkg := typ.PkgPath(); pkg != "" {
		return pkg + "." + typ.Name()
	} else {
		return typ.Name()
	}
}

type GeneralInterfaceFactory struct {
	interfaceType reflect.Type
	typeKey       string
	types         map[string]reflect.Type
	initializer   func(interface{}) error
}

func NewGeneralInterfaceFactory(interfaceType reflect.Type, typeKey string, initializer func(interface{}) error) *GeneralInterfaceFactory {
	return &GeneralInterfaceFactory{
		interfaceType: interfaceType,
		typeKey:       typeKey,
		types:         make(map[string]reflect.Type),
		initializer:   initializer,
	}
}

func (factory *GeneralInterfaceFactory) RegisterType(name string, instanceType reflect.Type) {
	factory.types[name] = instanceType
}

func (factory *GeneralInterfaceFactory) GetInstanceType() reflect.Type {
	return factory.interfaceType
}

func (factory *GeneralInterfaceFactory) Create(data map[string]interface{}) (interface{}, error) {
	var instance interface{}
	if typeName, ok := data[factory.typeKey].(string); !ok || typeName == "" {
		return nil, fmt.Errorf("missing type key: key=%q, map=%v", factory.typeKey, data)
	} else if instanceType, found := factory.types[typeName]; !found {
		return nil, fmt.Errorf("unknown type: %q", typeName)
	} else if instanceType.Kind() == reflect.Ptr {
		instance = reflect.New(instanceType.Elem()).Interface()
	} else {
		instance = reflect.New(instanceType).Interface()
	}
	if err := UnmarshalMap(instance, data); err != nil {
		return nil, fmt.Errorf("unmarshal map fail: %s", err.Error())
	}
	if factory.initializer != nil {
		if err := factory.initializer(instance); err != nil {
			return nil, fmt.Errorf("initialize fail: %s", err.Error())
		}
	}
	return instance, nil
}
