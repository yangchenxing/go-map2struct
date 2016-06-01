package map2struct

import (
	"encoding"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Second)

	percentageFloatPattern = regexp.MustCompile("^[0-9]+(\\.[0-9]+)?%$")

	timeLayouts = []string{
		"2006-01-02:15:04:05",
		"2006-01-02:15:04:05-0700",
		time.ANSIC,       // "Mon Jan _2 15:04:05 2006"
		time.UnixDate,    // "Mon Jan _2 15:04:05 MST 2006"
		time.RubyDate,    // "Mon Jan 02 15:04:05 -0700 2006"
		time.RFC822,      // "02 Jan 06 15:04 MST"
		time.RFC822Z,     // "02 Jan 06 15:04 -0700"
		time.RFC850,      // "Monday, 02-Jan-06 15:04:05 MST"
		time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
		time.Kitchen,     // "3:04PM"
		time.Stamp,       // "Jan _2 15:04:05"
		time.StampMilli,  // "Jan _2 15:04:05.000"
		time.StampMicro,  // "Jan _2 15:04:05.000000"
		time.StampNano,   // "Jan _2 15:04:05.000000000"
	}

	valueTrue  = reflect.ValueOf(true)
	valueFalse = reflect.ValueOf(false)
)

// Unmarshal unmarshal map[string]interface{} to a struct instance
func Unmarshal(dest, src interface{}) error {
	return unmarshal(rvalue(dest), reflect.ValueOf(src))
}

func unmarshal(dest, src reflect.Value) error {
	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}
	switch dest.Type() {
	case timeType:
		return unmarshalTime(dest, src)
	case durationType:
		return unmarshalDuration(dest, src)
	}
	if dest.CanAddr() {
		if textUnmarshaler, ok := dest.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return unmarshalText(textUnmarshaler, src)
		}
	}
	var unmarshalMethod func(reflect.Value, reflect.Value) error
	switch dest.Kind() {
	case reflect.Bool:
		unmarshalMethod = unmarshalBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		unmarshalMethod = unmarshalInt
	case reflect.Float32, reflect.Float64:
		unmarshalMethod = unmarshalFloat
	case reflect.Array:
		unmarshalMethod = unmarshalArray
	case reflect.Map:
		unmarshalMethod = unmarshalMap
	case reflect.Slice:
		unmarshalMethod = unmarshalSlice
	case reflect.String:
		unmarshalMethod = unmarshalString
	case reflect.Struct:
		unmarshalMethod = unmarshalStruct
	case reflect.Ptr:
		unmarshalMethod = unmarshalPtr
	case reflect.Interface:
		unmarshalMethod = unmarshalInterface
	}
	if unmarshalMethod == nil {
		return fmt.Errorf("unsupported kind: %s", dest.Kind())
	}
	return unmarshalMethod(dest, src)
}

func unmarshalBool(dest, src reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		dest.SetBool(src.Bool())
		return nil
	case reflect.String:
		boolean, err := strconv.ParseBool(src.String())
		if err != nil {
			return fmt.Errorf("invalid bool text: %s", src.String())
		}
		dest.SetBool(boolean)
		return nil
	}
	return badtype("bool/string", src)
}

func unmarshalInt(dest, src reflect.Value) error {
	srcKind := src.Kind()
	destKind := dest.Kind()
	if destKind >= reflect.Int && destKind <= reflect.Int64 {
		switch {
		case srcKind >= reflect.Int && srcKind <= reflect.Int64:
			dest.SetInt(src.Int())
		case srcKind >= reflect.Uint && srcKind <= reflect.Uint64:
			dest.SetInt(int64(src.Uint()))
		case srcKind == reflect.Float32 || srcKind == reflect.Float64:
			dest.SetInt(int64(src.Float()))
		case srcKind == reflect.String:
			intValue, err := parseIntText(src.String())
			if err != nil {
				return err
			}
			dest.SetInt(intValue)
		default:
			return badtype("int/string", src)
		}
	} else {
		switch {
		case srcKind >= reflect.Int && srcKind <= reflect.Int64:
			dest.SetUint(uint64(src.Int()))
		case srcKind >= reflect.Uint && srcKind <= reflect.Uint64:
			dest.SetUint(src.Uint())
		case srcKind == reflect.Float32 || srcKind == reflect.Float64:
			dest.SetUint(uint64(src.Float()))
		case srcKind == reflect.String:
			intValue, err := parseUintText(src.String())
			if err != nil {
				return err
			}
			dest.SetUint(uint64(intValue))
		default:
			return badtype("int/string", src)
		}
	}
	return nil
}

func unmarshalFloat(dest, src reflect.Value) error {
	srcKind := src.Kind()
	switch {
	case srcKind >= reflect.Int && srcKind <= reflect.Int64:
		dest.SetFloat(float64(src.Int()))
	case srcKind >= reflect.Uint && srcKind <= reflect.Uint64:
		dest.SetFloat(float64(src.Uint()))
	case srcKind == reflect.Float32 || srcKind == reflect.Float64:
		dest.SetFloat(src.Float())
	case srcKind == reflect.String:
		text := src.String()
		if text == "Inf" || text == "+Inf" {
			dest.SetFloat(math.Inf(1))
		} else if text == "-Inf" {
			dest.SetFloat(math.Inf(-1))
		} else if text == "NaN" {
			dest.SetFloat(math.NaN())
		} else if percentageFloatPattern.MatchString(text) {
			floatValue, _ := strconv.ParseFloat(text[:len(text)-1], 64)
			dest.SetFloat(floatValue / 100)
		} else {
			floatValue, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return err
			}
			dest.SetFloat(floatValue)
		}
	default:
		return badtype("int/float/string", src)
	}
	return nil
}

func unmarshalArray(dest, src reflect.Value) error {
	srcKind := src.Kind()
	if srcKind != reflect.Slice && srcKind != reflect.Array {
		return badtype("array/slice", src)
	} else if src.Len() != dest.Len() {
		return fmt.Errorf("array length mismatch: %d vs. %d", src.Len(), dest.Len())
	}
	return copySlice(dest, src)
}

func unmarshalInterface(dest, src reflect.Value) error {
	// interface{}
	if dest.Type().NumMethod() == 0 {
		dest.Set(src)
		return nil
	}
	// nil
	if !src.IsValid() {
		dest.Set(reflect.Zero(dest.Type()))
		return nil
	}
	// 其他直接赋值情况
	if dest.Type() == src.Type() || src.Type().Implements(dest.Type()) {
		dest.Set(src)
		return nil
	}
	// 非直接赋值情况
	if data, ok := src.Interface().(map[string]interface{}); !ok {
		return badtype("map[string]interface{}", src)
	} else if instance, err := createByFactory(dest.Type(), data); err != nil {
		return err
	} else {
		dest.Set(reflect.ValueOf(instance))
	}
	return nil
}

func unmarshalMap(dest, src reflect.Value) error {
	if src.Kind() == reflect.Slice && dest.Type().Elem().Kind() == reflect.Bool {
		return unmarshalSet(dest, src)
	}
	if src.Kind() != reflect.Map {
		return badtype("map", src)
	}
	if dest.IsNil() {
		dest.Set(reflect.MakeMap(dest.Type()))
	}
	keyType := dest.Type().Key()
	valueType := dest.Type().Elem()
	for _, srcKey := range src.MapKeys() {
		destKey := reflect.New(keyType).Elem()
		if err := unmarshal(destKey, srcKey); err != nil {
			return fmt.Errorf("unmarshal map index [%s] key error: %s",
				srcKey.Interface(), err.Error())
		}
		destValue := reflect.New(valueType).Elem()
		if err := unmarshal(destValue, src.MapIndex(srcKey)); err != nil {
			return fmt.Errorf("unmarshal map index [%s] value error: %s",
				srcKey.Interface(), err.Error())
		}
		dest.SetMapIndex(destKey, destValue)
	}
	return nil
}

func unmarshalSet(dest, src reflect.Value) error {
	if dest.IsNil() {
		dest.Set(reflect.MakeMap(dest.Type()))
	}
	keyType := dest.Type().Key()
	for i := 0; i < src.Len(); i++ {
		destKey := reflect.New(keyType).Elem()
		if err := unmarshal(destKey, src.Index(i)); err != nil {
			return fmt.Errorf("unmarshal set item [%d] error: %s", i, err.Error())
		}
		dest.SetMapIndex(destKey, valueTrue)
	}
	return nil
}

func unmarshalSlice(dest, src reflect.Value) error {
	srcKind := src.Kind()
	if srcKind != reflect.Slice && srcKind != reflect.Array {
		return badtype("array/slice", src)
	}
	if dest.Len() < src.Len() {
		dest.Set(reflect.MakeSlice(dest.Type(), src.Len(), src.Len()))
	}
	return copySlice(dest, src)
}

func unmarshalString(dest, src reflect.Value) error {
	if src.Kind() != reflect.String {
		return badtype("string", src)
	}
	dest.SetString(src.String())
	return nil
}

func unmarshalStruct(dest, src reflect.Value) error {
	if dest.Type() == src.Type() {
		dest.Set(src)
		return nil
	}
	data, ok := src.Interface().(map[string]interface{})
	if !ok {
		return badtype("map[string]interface{}", src)
	}

	typ := dest.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			if err := unmarshal(dest.Field(i), src); err != nil {
				return fmt.Errorf("unmarshal anonymous field %q fail: type=%q, error=%q",
					field.Name, typ, err.Error())
			}
			continue
		}
		value, found := data[field.Name]
		if !found {
			continue
		}
		switch field.Type.Kind() {
		case reflect.Ptr:
			if dest.Field(i).IsNil() {
				dest.Field(i).Set(reflect.New(field.Type.Elem()))
			}
			if err := unmarshal(dest.Field(i).Elem(), reflect.ValueOf(value)); err != nil {
				return fmt.Errorf("unmarshal field %s fail: %s", field.Name, err.Error())
			}
		case reflect.Interface:
			if err := unmarshalInterface(dest.Field(i), reflect.ValueOf(value)); err != nil {
				return fmt.Errorf("unmarshal field %s fail: %s", field.Name, err.Error())
			}
		default:
			if err := unmarshal(dest.Field(i), reflect.ValueOf(value)); err != nil {
				return fmt.Errorf("unmarshal field %s fail: %s", field.Name, err.Error())
			}
		}
	}
	return nil
}

func unmarshalPtr(dest, src reflect.Value) error {
	// always non-nil
	// if dest.IsNil() {
	// 	dest.Set(reflect.New(dest.Type().Elem()))
	// }
	return unmarshal(reflect.Indirect(dest), src)
}

func unmarshalText(dest encoding.TextUnmarshaler, src reflect.Value) error {
	if src.Kind() == reflect.String {
		return dest.UnmarshalText([]byte(src.String()))
	} else if bytes, ok := src.Interface().([]byte); ok {
		return dest.UnmarshalText(bytes)
	}
	return badtype("string/[]byte", src)
}

func unmarshalTime(dest, src reflect.Value) error {
	if src.Kind() != reflect.String {
		return badtype("string", src)
	}
	text := src.String()
	for _, layout := range timeLayouts {
		if len(layout) == len(text) {
			if timeValue, err := time.Parse(layout, text); err == nil {
				dest.Set(reflect.ValueOf(timeValue))
				return nil
			}
		}
	}
	return fmt.Errorf("unknown time layout: %s", text)
}

func unmarshalDuration(dest, src reflect.Value) error {
	if src.Kind() != reflect.String {
		return badtype("string", src)
	}
	text := src.String()
	if text == "genesis" {
		dest.SetInt(math.MinInt64)
		return nil
	} else if text == "doomsday" {
		dest.SetInt(math.MaxInt64)
		return nil
	} else if durationValue, err := time.ParseDuration(text); err == nil {
		dest.Set(reflect.ValueOf(durationValue))
		return nil
	}
	return fmt.Errorf("invalid duration: %s", text)
}

func parseIntText(text string) (int64, error) {
	if text == "0" {
		return 0, nil
	} else if strings.HasPrefix(text, "0x") {
		return strconv.ParseInt(text[2:], 16, 64)
	} else if strings.HasPrefix(text, "0") {
		return strconv.ParseInt(text[1:], 8, 64)
	}
	return strconv.ParseInt(text, 10, 64)
}

func parseUintText(text string) (uint64, error) {
	if text == "0" {
		return 0, nil
	} else if strings.HasPrefix(text, "0x") {
		return strconv.ParseUint(text[2:], 16, 64)
	} else if strings.HasPrefix(text, "0") {
		return strconv.ParseUint(text[1:], 8, 64)
	}
	return strconv.ParseUint(text, 10, 64)
}

func copySlice(dest, src reflect.Value) error {
	for i, len := 0, dest.Len(); i < len; i++ {
		if err := unmarshal(dest.Index(i), src.Index(i)); err != nil {
			return fmt.Errorf("copy index [%d] error: %s", i, err.Error())
		}
	}
	return nil
}

func rvalue(value interface{}) reflect.Value {
	return indirect(reflect.ValueOf(value))
}

func indirect(value reflect.Value) reflect.Value {
	if value.Kind() != reflect.Ptr {
		if value.CanAddr() {
			pv := value.Addr()
			if _, ok := pv.Interface().(encoding.TextUnmarshaler); ok {
				return pv
			}
		}
		return value
	}
	if value.IsNil() {
		value.Set(reflect.New(value.Type().Elem()))
	}
	return indirect(reflect.Indirect(value))
}

func badtype(expected string, value reflect.Value) error {
	return fmt.Errorf("expect %s but found %q(%s)",
		expected, value.Type().Name(), value.Kind())
}
