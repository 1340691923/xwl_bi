package parser

import (
	"fmt"
	"github.com/1340691923/xwl_bi/engine/logs"
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
	timeZone     *time.Location
	pool         sync.Pool
}

// NewParserPool creates a parser pool
func NewParserPool() (pp *Pool, err error) {
	var tz =  time.Local

	pp = &Pool{
		timeZone:  tz,
	}

	return
}

func ParseKafkaData(pool *Pool, data []byte) (metric *FastjsonMetric, err error) {
	jsonParser := pool.Get()
	defer pool.Put(jsonParser)
	metric, err = jsonParser.Parse(data)
	return
}

// Get returns a Parser from pp.
//
// The Parser must be Put to pp after use.
func (pp *Pool) Get() Parser {
	v := pp.pool.Get()
	if v == nil {
		return &FastjsonParser{pp: pp}
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

func (pp *Pool) ParseDateTime(key string, val string) (t time.Time, err error) {

	var t2 time.Time
	if val == "" {
		err = ErrParseDateTime
		return
	}


	if t2, err = time.ParseInLocation("2006-01-02 15:04:05", val, pp.timeZone); err != nil {
		err = ErrParseDateTime
		return
	}
	t = t2.UTC()
	return
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

func parseInLocation(val string, loc *time.Location) (t time.Time, layout string) {
	var err error
	var lay string
	for _, lay = range Layouts {
		if t, err = time.ParseInLocation(lay, val, loc); err == nil {
			t = t.UTC()
			layout = lay
			return
		}
	}
	return
}

func UnixFloat(sec float64) (t time.Time) {
	if sec < 0 || sec >= 4294967296.0 {
		return Epoch
	}
	i, f := math.Modf(sec)
	return time.Unix(int64(i), int64(f*1e9)).UTC()
}
