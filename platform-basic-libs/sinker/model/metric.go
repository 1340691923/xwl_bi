package model

import (
	"sync"
)

type Metric interface {
	GetInt(key string, nullable bool) (val interface{})
	GetFloat(key string, nullable bool) (val interface{})
	GetString(key string, nullable bool) (val interface{})
	GetDateTime(key string, nullable bool) (val interface{})
	GetElasticDateTime(key string, nullable bool) (val interface{})
	GetArray(key string, t int) (val interface{})
	GetNewKeys(knownKeys *sync.Map, newKeys *sync.Map) bool
}

type DimMetrics struct {
	Dims   []*ColumnWithType
	Fields []*ColumnWithType
}

type ColumnWithType struct {
	Name       string
	Type       int
	Nullable   bool
	SourceName string
}
