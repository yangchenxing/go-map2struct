package map2struct

import (
	"fmt"
	"reflect"
)

// Factory represents a instance factory used in unmarshaling.
type Factory interface {
	// GetInstanceType returns the type of the instance create by this factory.
	// The instance type is used to index the factory.
	GetInstanceType() reflect.Type

	// Create returns the new instance.
	Create(map[string]interface{}) (interface{}, error)
}

var (
	factories = make(map[string]Factory)
)

// RegisterFactory register factories.
func RegisterFactory(factory Factory) {
	factories[getTypeName(factory.GetInstanceType())] = factory
}

func createByFactory(typ reflect.Type, data map[string]interface{}) (interface{}, error) {
	typeName := getTypeName(typ)
	if factory := factories[typeName]; factory != nil {
		return factory.Create(data)
	}
	return nil, fmt.Errorf("unregistered type: %q", typeName)
}

func getTypeName(typ reflect.Type) string {
	if pkg := typ.PkgPath(); pkg != "" {
		return pkg + "." + typ.Name()
	}
	return typ.Name()
}

// GeneralInterfaceFactory provides a general factory for interface.
// For creating instance, the input map need a key to specified the type implelement the interface.
// The factory new a instance with the type associated the type key, and unmarshal the map to the instance struct.
type GeneralInterfaceFactory struct {
	interfaceType reflect.Type
	typeKey       string
	types         map[string]reflect.Type
	initializer   func(interface{}) error
}

// NewGeneralInterfaceFactory creates a GeneralInterfaceFactory instance.
func NewGeneralInterfaceFactory(interfaceType reflect.Type, typeKey string, initializer func(interface{}) error) *GeneralInterfaceFactory {
	return &GeneralInterfaceFactory{
		interfaceType: interfaceType,
		typeKey:       typeKey,
		types:         make(map[string]reflect.Type),
		initializer:   initializer,
	}
}

// RegisterType register new type and its associated key.
func (factory *GeneralInterfaceFactory) RegisterType(name string, instanceType reflect.Type) {
	factory.types[name] = instanceType
}

// GetInstanceType returns the interface type.
func (factory *GeneralInterfaceFactory) GetInstanceType() reflect.Type {
	return factory.interfaceType
}

// Create creates a new instance implement the interface.
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
	if err := Unmarshal(instance, data); err != nil {
		return nil, fmt.Errorf("unmarshal map fail: %s", err.Error())
	}
	if factory.initializer != nil {
		if err := factory.initializer(instance); err != nil {
			return nil, fmt.Errorf("initialize fail: %s", err.Error())
		}
	}
	return instance, nil
}
