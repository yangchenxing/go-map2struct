package map2struct

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

type TestEmbedTypeB struct {
	I int
}

type TestTypeA struct {
	IntA      int
	IntB      uint32
	IntC      int64
	IntHex    uint32
	IntOct    uint16
	FloatA    float32
	BoolA     bool
	BoolB     bool
	ArrayA    [2]int
	SliceA    []int
	MapA      map[string]string
	DurationA time.Duration
	TimeA     time.Time
	EmbedA    TestEmbedTypeB
	EmbedB    *TestEmbedTypeB
}

func (a TestTypeA) Equal(b TestTypeA) bool {
	result := a.IntA == b.IntA &&
		a.IntB == b.IntB &&
		a.IntC == b.IntC &&
		a.IntHex == b.IntHex &&
		a.IntOct == b.IntOct &&
		a.FloatA == b.FloatA &&
		a.BoolA == b.BoolA &&
		a.BoolB == b.BoolB &&
		a.ArrayA == b.ArrayA &&
		a.DurationA == b.DurationA &&
		a.TimeA == b.TimeA &&
		a.EmbedA == b.EmbedA &&
		len(a.MapA) == len(b.MapA) &&
		len(a.SliceA) == len(b.SliceA) &&
		(a.EmbedB == nil && b.EmbedB == nil || a.EmbedB != nil && b.EmbedB != nil)
	if !result {
		return false
	}
	for key, value := range a.MapA {
		if value != b.MapA[key] {
			return false
		}
	}
	for i, value := range a.SliceA {
		if value != b.SliceA[i] {
			return false
		}
	}
	if a.EmbedB != nil && *a.EmbedB != *b.EmbedB {
		return false
	}
	return true
}

func TestUnmarshal(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	src := map[string]interface{}{
		"IntA":   1,
		"IntB":   uint32(2),
		"IntC":   "3",
		"IntHex": "0x0fffffff",
		"IntOct": "0755",
		"FloatA": 1.0,
		"BoolA":  false,
		"BoolB":  "true",
		"ArrayA": []interface{}{1, 2},
		"SliceA": []interface{}{uint(3), uint(4)},
		"MapA": map[string]interface{}{
			"Key": "Value",
		},
		"DurationA": "1h",
		"TimeA":     now.Format("2006-01-02:15:04:05-0700"),
		"EmbedA": map[string]interface{}{
			"I": "99",
		},
		"EmbedB": map[string]interface{}{
			"I": "98",
		},
	}
	// t.Log(src)
	expect := TestTypeA{
		IntA:   1,
		IntB:   2,
		IntC:   3,
		IntHex: 0x0fffffff,
		IntOct: 0755,
		FloatA: 1.0,
		BoolA:  false,
		BoolB:  true,
		ArrayA: [2]int{1, 2},
		SliceA: []int{3, 4},
		MapA: map[string]string{
			"Key": "Value",
		},
		DurationA: time.Hour,
		TimeA:     now,
		EmbedA: TestEmbedTypeB{
			I: 99,
		},
		EmbedB: &TestEmbedTypeB{
			I: 98,
		},
	}
	var actual *TestTypeA
	if err := Unmarshal(&actual, src); err != nil {
		t.Errorf("unmarshal map fail: %s\n", err.Error())
		return
	}
	if !actual.Equal(expect) {
		t.Errorf("unmarshal result fail: expect=%s, actual=%s\n", jsonify(expect), jsonify(actual))
		return
	}
}

func TestUnmarshalBool(t *testing.T) {
	var a bool
	// invalid string value
	if err := Unmarshal(&a, "foobar"); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	// invalid type
	a = false
	if err := Unmarshal(&a, 1); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
}

func TestUnmarshalInt(t *testing.T) {
	var a int
	// float
	if err := Unmarshal(&a, 2.0); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a != 2 {
		t.Error("unexpected unmarshal result:", a)
		return
	}
	// invalid string value
	a = 0
	if err := Unmarshal(&a, "abc"); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	// invalid type
	a = 0
	if err := Unmarshal(&a, true); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}

	if err := Unmarshal(&a, true); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	var b uint
	// int
	if err := Unmarshal(&b, 3); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if b != 3 {
		t.Error("unexpected unmarshal result:", b)
		return
	}
	// float
	b = 0
	if err := Unmarshal(&b, 2.0); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if b != 2 {
		t.Error("unexpected unmarshal result:", b)
		return
	}
	// invalid string value
	b = 0
	if err := Unmarshal(&b, "abc"); err == nil {
		t.Error("unexpected unmarshal success:", b)
		return
	}
	// invalid type
	b = 0
	if err := Unmarshal(&b, true); err == nil {
		t.Error("unexpected unmarshal success:", b)
		return
	}
}

func TestUnmarshalFloat(t *testing.T) {
	var a float64
	// int
	if err := Unmarshal(&a, 10); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a != 10 {
		t.Error("unexpected unmarshal resutl:", a)
		return
	}
	// uint
	a = 0
	if err := Unmarshal(&a, uint(20)); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a != 20 {
		t.Error("unexpected unmarshal resutl:", a)
		return
	}
	// string
	a = 0
	if err := Unmarshal(&a, "2.1"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a != 2.1 {
		t.Error("unexpected unmarshal resutl:", a)
		return
	}
	// NaN
	a = 0
	if err := Unmarshal(&a, "NaN"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if !math.IsNaN(a) {
		t.Error("unexpected unmarshal resutl:", a)
		return
	}
	// Percentage
	a = 0
	if err := Unmarshal(&a, "15%"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a != 0.15 {
		t.Error("unexpected unmarshal resutl:", a)
		return
	}
	// Inf
	a = 0
	if err := Unmarshal(&a, "+Inf"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if !math.IsInf(a, 1) {
		t.Error("unexpected unmarshal resutl:", a)
		return
	}
	a = 0
	if err := Unmarshal(&a, "-Inf"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if !math.IsInf(a, -1) {
		t.Error("unexpected unmarshal resutl:", a)
		return
	}
	// invalid string
	a = 0
	if err := Unmarshal(&a, "abc"); err == nil {
		t.Error("unexpected unmarshal succes:", a)
		return
	}
	// badtype
	a = 0
	if err := Unmarshal(&a, true); err == nil {
		t.Error("unexpected unmarshal succes:", a)
		return
	}
}

func TestUnmarshalArray(t *testing.T) {
	var a [2]int
	// badtype
	if err := Unmarshal(&a, 1); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	// unmatched length
	if err := Unmarshal(&a, []int{1, 2, 3}); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
}

func TestUnmarshalInterface(t *testing.T) {
	var a stringer
	var c stringer
	// nil
	if err := Unmarshal(&a, c); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	}
	// interface{}
	var b interface{}
	if err := Unmarshal(&b, "string"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	}
	// same type
	c = foo{Text: "hello"}
	if err := Unmarshal(&a, c); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	}
	// implements
	d := foo{Text: "world"}
	if err := Unmarshal(&a, d); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	}
	// invalid type
	a = nil
	if err := Unmarshal(&a, 1); err == nil {
		t.Error("unexpected unmarshal success", a)
		return
	}
}

func TestUnmarshalMap(t *testing.T) {
	// set
	var a map[string]bool
	if err := Unmarshal(&a, []string{"a", "b"}); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if len(a) != 2 || !a["a"] || !a["b"] {
		t.Error("unexpected unmarshal result:", a)
		return
	}
	a = nil
	if err := Unmarshal(&a, []int{1, 2}); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	// badtype
	a = nil
	if err := Unmarshal(&a, ""); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	// invalid key type
	a = nil
	if err := Unmarshal(&a, map[int]int{1: 2}); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	// invalid value type
	a = nil
	if err := Unmarshal(&a, map[string]int{"1": 2}); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
}

func TestUnmarshalSlice(t *testing.T) {
	// badtype
	var a []int
	if err := Unmarshal(&a, ""); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
}

type S1 struct {
	Duration time.Duration
}

type S2 struct {
	S1
}

type S3 struct {
	S1 *string
}

type S4 struct {
	S stringer
}

type S5 struct {
	S1 *S1
}

func TestUnmarshalStruct(t *testing.T) {
	// same type
	var a, b foo
	b.Text = "hello"
	if err := Unmarshal(&a, b); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a.Text != "hello" {
		t.Error("unexpected unmarshal result:", a)
		return
	}
	// badtype
	a.Text = ""
	if err := Unmarshal(&a, 1); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}

	var src map[string]interface{}
	var c S2
	// anonymous success
	src = map[string]interface{}{
		"Duration": "1s",
	}
	if err := Unmarshal(&c, src); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if c.Duration != time.Second {
		t.Error("unexpected unmarshal result:", c)
		return
	}
	// anonymous fail
	src = map[string]interface{}{
		"Duration": true,
	}
	if err := Unmarshal(&c, src); err == nil {
		t.Error("unexpected unmarshal success:", c)
		return
	}
	// missing field
	c.Duration = 0
	if err := Unmarshal(&c, map[string]interface{}{}); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if c.Duration != 0 {
		t.Error("unexpected unmarshal result:", c)
		return
	}
	// ptr fail
	src = map[string]interface{}{
		"S1": true,
	}
	var d S3
	if err := Unmarshal(&d, src); err == nil {
		t.Error("unexpected unmarshal success:", c)
		return
	}
	// interface fail
	src = map[string]interface{}{
		"S": true,
	}
	var e S4
	if err := Unmarshal(&e, src); err == nil {
		t.Error("unexpected unmarshal success:", c)
		return
	}
}

type UT struct {
	text string
}

func TestUnmarshalPtr(t *testing.T) {
	var a S5
	src := map[string]interface{}{
		"Duration": "2s",
	}
	if err := Unmarshal(&a, src); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	}

	var b uintptr
	if err := Unmarshal(&b, nil); err == nil {
		t.Error("unexpected unmarshal success")
		return
	} else if err.Error() != "unsupported kind: uintptr" {
		t.Error("unexpected error:", err.Error())
		return
	}
}

func (tm *UT) UnmarshalText(content []byte) error {
	tm.text = string(content)
	return nil
}

func TestUnmarshalText(t *testing.T) {
	var a UT
	if err := Unmarshal(&a, "hello"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a.text != "hello" {
		t.Error("unexpected unmarshal result:", a)
		return
	}
	a.text = ""
	if err := Unmarshal(&a, []byte("world")); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a.text != "world" {
		t.Error("unexpected unmarshal result:", a)
		return
	}
	a.text = ""
	if err := Unmarshal(&a, 1); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
}

func TestUnmarshalTime(t *testing.T) {
	var a time.Time
	if err := Unmarshal(&a, true); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
	if err := Unmarshal(&a, "abc"); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
}

func TestUnmarshalDuration(t *testing.T) {
	var a time.Duration
	if err := Unmarshal(&a, "genesis"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a != time.Duration(math.MinInt64) {
		t.Error("unexpected unmarshal result:", a)
		return
	}
	if err := Unmarshal(&a, "doomsday"); err != nil {
		t.Error("unmarshal fail:", err.Error())
		return
	} else if a != time.Duration(math.MaxInt64) {
		t.Error("unexpected unmarshal result:", a)
		return
	}
}

func TestParseInt(t *testing.T) {
	if v, err := parseIntText("0"); err != nil {
		t.Error("parseIntText fail:", err.Error())
		return
	} else if v != 0 {
		t.Error("unexpected parseIntText result:", v)
		return
	}
	if v, err := parseIntText("0x10"); err != nil {
		t.Error("parseIntText fail:", err.Error())
		return
	} else if v != 16 {
		t.Error("unexpected parseIntText result:", v)
		return
	}
	if v, err := parseIntText("010"); err != nil {
		t.Error("parseIntText fail:", err.Error())
		return
	} else if v != 8 {
		t.Error("unexpected parseIntText result:", v)
		return
	}
}

func TestParseUint(t *testing.T) {
	if v, err := parseUintText("0"); err != nil {
		t.Error("parseUintText fail:", err.Error())
		return
	} else if v != 0 {
		t.Error("unexpected parseUintText result:", v)
		return
	}
	if v, err := parseUintText("0x10"); err != nil {
		t.Error("parseUintText fail:", err.Error())
		return
	} else if v != 16 {
		t.Error("unexpected parseUintText result:", v)
		return
	}
	if v, err := parseUintText("010"); err != nil {
		t.Error("parseUintText fail:", err.Error())
		return
	} else if v != 8 {
		t.Error("unexpected parseUintText result:", v)
		return
	}
}

func TestCopySlice(t *testing.T) {
	a := make([]string, 10)
	b := make([]int, 10)
	if err := Unmarshal(&a, b); err == nil {
		t.Error("unexpected unmarshal success:", a)
		return
	}
}

func jsonify(i interface{}) string {
	c, _ := json.Marshal(i)
	return string(c)
}
