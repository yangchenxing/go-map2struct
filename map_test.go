package map2struct

import (
	"encoding/json"
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

func TestUnmarshalMap(t *testing.T) {
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
	if err := UnmarshalMap(&actual, src); err != nil {
		t.Errorf("unmarshal map fail: %s\n", err.Error())
		return
	}
	if !actual.Equal(expect) {
		t.Errorf("unmarshal result fail: expect=%s, actual=%s\n", jsonify(expect), jsonify(actual))
		return
	}
}

func jsonify(i interface{}) string {
	c, _ := json.Marshal(i)
	return string(c)
}
