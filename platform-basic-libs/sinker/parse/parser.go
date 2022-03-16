package parser

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"math"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	Layouts = []string{
		"2006-01-02 15:04:05",
	}
	Epoch            = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	ErrParseDateTime = errors.Errorf("value doesn't contain DateTime")
)

// Parse is the Parser interface
type Parser interface {
	Parse(bs []byte) (metric *FastjsonMetric, err error)
}

// Pool may be used for pooling Parsers for similarly typed JSONs.
type Pool struct {
	name     string
	timeZone *time.Location
	pool     sync.Pool
}

// NewParserPool creates a parser pool
func NewParserPool(name string) (pp *Pool, err error) {

	pp = &Pool{
		name:     name,
		timeZone: time.Local,
	}

	return
}

// Get returns a Parser from pp.
//
// The Parser must be Put to pp after use.
func (pp *Pool) Get() Parser {
	v := pp.pool.Get()
	if v == nil {
		switch pp.name {
		case "fastjson":
			return &FastjsonParser{}
		default:
			return &FastjsonParser{}
		}
	}
	return v.(Parser)
}

// Put returns p to pp.
//
// p and objects recursively returned from p cannot be used after p
// is put into pp.
func (pp *Pool) Put(p Parser) {
	pp.pool.Put(p)
}

func makeArray(typ int) (val interface{}) {
	switch typ {
	case Int:
		val = []int64{}
	case Float:
		val = []float64{}
	case String:
		val = []string{}
	case DateTime:
		val = []time.Time{}
	default:
		logs.Logger.Sugar().Errorf(fmt.Sprintf("LOGIC ERROR: unsupported array type %v", typ))
	}
	return
}

func parseInLocation(val string, loc *time.Location) (t time.Time, err error) {

	if t, err = time.ParseInLocation(util.TimeFormat, val, loc); err == nil {
		t = t.UTC()
		return
	}

	return t, err
}

func UnixFloat(sec float64) (t time.Time) {
	if sec < 0 || sec >= 4294967296.0 {
		return Epoch
	}
	i, f := math.Modf(sec)
	return time.Unix(int64(i), int64(f*1e9)).UTC()
}
