# go-map2struct

[![Go Report Card](https://goreportcard.com/badge/github.com/yangchenxing/go-map2struct)](https://goreportcard.com/report/github.com/yangchenxing/go-map2struct)
[![Build Status](https://travis-ci.org/yangchenxing/go-map2struct.svg?branch=master)](https://travis-ci.org/yangchenxing/go-map2struct)
[![GoDoc](http://godoc.org/github.com/yangchenxing/go-map2struct?status.svg)](http://godoc.org/github.com/yangchenxing/go-map2struct)
[![Coverage Status](https://coveralls.io/repos/github/yangchenxing/go-map2struct/badge.svg?branch=master)](https://coveralls.io/github/yangchenxing/go-map2struct?branch=master)

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
