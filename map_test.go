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
	DurationA time.Duration
	TimeA     time.Time
	EmbedA    TestEmbedTypeB
}

func TestUnmarshalMap(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	src := map[string]interface{}{
		"IntA":      1,
		"IntB":      2,
		"IntC":      "3",
		"IntHex":    "0x0fffffff",
		"IntOct":    "0755",
		"FloatA":    1.0,
		"BoolA":     false,
		"BoolB":     "true",
		"DurationA": "1h",
		"TimeA":     now.Format("2006-01-02:15:04:05"),
		"EmbedA": map[string]interface{}{
			"I": "99",
		},
	}
	t.Log(src)
	expect := TestTypeA{
		IntA:      1,
		IntB:      2,
		IntC:      3,
		IntHex:    0x0fffffff,
		IntOct:    0755,
		FloatA:    1.0,
		BoolA:     false,
		BoolB:     true,
		DurationA: time.Hour,
		TimeA:     now,
		EmbedA: TestEmbedTypeB{
			I: 99,
		},
	}
	var actual TestTypeA
	if err := UnmarshalMap(&actual, src); err != nil {
		t.Errorf("unmarshal map fail: %s\n", err.Error())
	}
	if actual.IntA != 1 || actual.IntB != 2 || actual.IntC != 3 || actual.IntHex != 0x0fffffff || actual.IntOct != 0755 ||
		actual.FloatA != 1.0 || actual.BoolA != false || actual.BoolB != true {
		t.Errorf("unmarshal result fail: expect=%s, actual=%s\n", jsonify(expect), jsonify(actual))
	}
}

func jsonify(i interface{}) string {
	c, _ := json.Marshal(i)
	return string(c)
}
