package parser

import (
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
)

type FastjsonParser struct {
	fjp fastjson.Parser
}

func (p *FastjsonParser) Parse(bs []byte) (metric *FastjsonMetric, err error) {
	var value *fastjson.Value
	if value, err = p.fjp.ParseBytes(bs); err != nil {
		err = errors.Wrapf(err, "")
		return
	}
	metric = &FastjsonMetric{value: value}
	return
}

type FastjsonMetric struct {
	value *fastjson.Value
}

func (c *FastjsonMetric) GetString(key string, nullable bool) (val interface{}) {
	v := c.value.Get(key)
	if v == nil || v.Type() == fastjson.TypeNull {
		if nullable {
			return
		}
		val = ""
		return
	}
	switch v.Type() {
	case fastjson.TypeString:
		b, _ := v.StringBytes()
		val = util.Bytes2str(b)
	default:
		val = v.String()
	}
	return
}

func (c *FastjsonMetric) GetFloat(key string, nullable bool) (val interface{}) {
	v := c.value.Get(key)
	if !fjCompatibleFloat(v) {
		val = getDefaultFloat(nullable)
		return
	}
	if val2, err := v.Float64(); err != nil {
		val = getDefaultFloat(nullable)
	} else {
		val = val2
	}
	return
}

func (c *FastjsonMetric) GetInt(key string, nullable bool) (val interface{}) {
	v := c.value.Get(key)
	if !fjCompatibleInt(v) {
		val = getDefaultInt(nullable)
		return
	}
	switch v.Type() {
	case fastjson.TypeTrue:
		val = int64(1)
	case fastjson.TypeFalse:
		val = int64(0)
	default:
		if val2, err := v.Int64(); err != nil {
			val = getDefaultInt(nullable)
		} else {
			val = val2
		}
	}
	return
}

func (c *FastjsonMetric) GetDateTime(key string, nullable bool) (val interface{}) {
	v := c.value.Get(key)
	if !fjCompatibleDateTime(v) {
		val = getDefaultDateTime(nullable)
		return
	}
	var err error
	switch v.Type() {
	case fastjson.TypeNumber:
		var f float64
		if f, err = v.Float64(); err != nil {
			val = getDefaultDateTime(nullable)
			return
		}
		val = UnixFloat(f)
	case fastjson.TypeString:
		var b []byte
		if b, err = v.StringBytes(); err != nil || len(b) == 0 {
			val = getDefaultDateTime(nullable)
			return
		}
		if val, err = c.ParseDateTime(util.Bytes2str(b)); err != nil {
			val = getDefaultDateTime(nullable)
		}
	default:
		val = getDefaultDateTime(nullable)
	}
	return
}

func (c *FastjsonMetric) ParseDateTime(val string) (t time.Time, err error) {

	var t2 time.Time
	if val == "" {
		err = ErrParseDateTime
		return
	}
	if t2, err = time.ParseInLocation(util.TimeFormat, val, time.Local); err != nil {
		err = ErrParseDateTime
		return
	}
	t = t2.UTC()
	return
}

func (c *FastjsonMetric) GetElasticDateTime(key string, nullable bool) (val interface{}) {
	t := c.GetDateTime(key, nullable)
	if t != nil {
		val = t.(time.Time).Unix()
	}
	return
}

func (c *FastjsonMetric) GetArray(key string, typ int) (val interface{}) {
	v := c.value.Get(key)
	val = makeArray(typ)
	if v == nil || v.Type() != fastjson.TypeArray {
		return
	}
	array, _ := v.Array()
	switch typ {
	case Int:
		for _, e := range array {
			var v int64
			if e.Type() == fastjson.TypeTrue {
				v = 1
			} else {
				v, _ = e.Int64()
			}
			val = append(val.([]int64), v)
		}
	case Float:
		for _, e := range array {
			v, _ := e.Float64()
			val = append(val.([]float64), v)
		}
	case String:
		for _, e := range array {
			var s string
			switch e.Type() {
			case fastjson.TypeNull:
				s = ""
			case fastjson.TypeString:
				b, _ := e.StringBytes()
				s = util.Bytes2str(b)
			default:
				s = e.String()
			}
			val = append(val.([]string), s)
		}
	case DateTime:
		for _, e := range array {
			var t time.Time
			switch e.Type() {
			case fastjson.TypeNumber:
				if f, err := e.Float64(); err != nil {
					t = Epoch
				} else {
					t = UnixFloat(f)
				}
			case fastjson.TypeString:
				if b, err := e.StringBytes(); err != nil || len(b) == 0 {
					t = Epoch
				} else {
					var err error
					if t, err = c.ParseDateTime(util.Bytes2str(b)); err != nil {
						t = Epoch
					}
				}
			default:
				t = Epoch
			}
			val = append(val.([]time.Time), t)
		}
	default:
		logs.Logger.Error("LOGIC ERROR: unsupported array", zap.Int("typ", typ))
	}
	return
}

func (c *FastjsonMetric) GetNewKeys(knownKeys *sync.Map, newKeys *sync.Map) (foundNew bool) {
	var obj *fastjson.Object
	var err error
	if obj, err = c.value.Object(); err != nil {
		return
	}
	obj.Visit(func(key []byte, v *fastjson.Value) {
		strKey := util.Bytes2str(key)
		if _, loaded := knownKeys.LoadOrStore(strKey, nil); !loaded {
			if typ := FjDetectType(v); typ != TypeUnknown {
				newKeys.Store(strKey, typ)
				foundNew = true
			} else {
				logs.Logger.Error("FastjsonMetric.GetNewKeys failed to detect field type", zap.String("key", strKey), zap.Reflect("value", v))
			}
		}
	})
	return
}

func (c *FastjsonMetric) GetParseObject() (v *fastjson.Value) {
	return c.value
}

func fjCompatibleInt(v *fastjson.Value) (ok bool) {
	if v == nil {
		return
	}
	switch v.Type() {
	case fastjson.TypeTrue:
		ok = true
	case fastjson.TypeFalse:
		ok = true
	case fastjson.TypeNumber:
		ok = true
	}
	return
}

func fjCompatibleFloat(v *fastjson.Value) (ok bool) {
	if v == nil {
		return
	}
	switch v.Type() {
	case fastjson.TypeNumber:
		ok = true
	}
	return
}

func fjCompatibleDateTime(v *fastjson.Value) (ok bool) {
	if v == nil {
		return
	}
	switch v.Type() {
	case fastjson.TypeNumber:
		ok = true
	case fastjson.TypeString:
		ok = true
	}
	return
}

func getDefaultInt(nullable bool) (val interface{}) {
	if nullable {
		return
	}
	val = int64(0)
	return
}

func getDefaultFloat(nullable bool) (val interface{}) {
	if nullable {
		return
	}
	val = float64(0.0)
	return
}

func getDefaultDateTime(nullable bool) (val interface{}) {
	if nullable {
		return
	}
	val = Epoch
	return
}

func FjDetectType(v *fastjson.Value) (typ int) {
	switch v.Type() {
	case fastjson.TypeNull:
	case fastjson.TypeTrue:
		typ = Int
	case fastjson.TypeFalse:
		typ = Int
	case fastjson.TypeNumber:
		typ = Float
		if _, err := v.Int64(); err == nil {
			typ = Int
		}
	case fastjson.TypeString:
		typ = String
		if val, err := v.StringBytes(); err == nil {
			if _, err := parseInLocation(util.Bytes2str(val), time.Local); err == nil {
				typ = DateTime
			}
		}
	case fastjson.TypeArray:
		if arr, err := v.Array(); err == nil && len(arr) > 0 {
			typ2 := FjDetectType(arr[0])
			switch typ2 {
			case Int:
				typ = IntArray
			case Float:
				typ = FloatArray
			case String:
				typ = StringArray
			case DateTime:
				typ = DateTimeArray
			}
		}
	default:
		typ = String
	}
	return
}
