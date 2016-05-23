# go-map2struct

go-map2struct convert map[string]interface{} to struct.

## Example

    type Foo struct {
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
      EmbedA    Bar
    }
    type Bar struct {
      I int
    }
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
      "TimeA":     now.Format("2006-01-02:15:04:05-0700"),
      "EmbedA": map[string]interface{}{
        "I": "99",
      },
    }
    var dest Foo
    UnmarshalMap(dest, src)
