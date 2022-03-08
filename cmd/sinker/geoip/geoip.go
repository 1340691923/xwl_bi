package geoip

import (
	_ "embed"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/oschwald/maxminddb-golang"
	"go.uber.org/zap"
	"net"
	"sync"
)

//go:embed GeoLite2-City.mmdb
var GeoipMmdbByte []byte

type Res struct {
	City struct {
		Names struct {
			ZhCN string `maxminddb:"zh-CN"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		Names struct {
			ZhCN string `maxminddb:"zh-CN"`
		} `maxminddb:"names"`
	} `maxminddb:"country"`
	Subdivisions []struct {
		Names struct {
			ZhCN string `maxminddb:"zh-CN"`
		} `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
}

func (this *Res) reset() {
	this.City.Names.ZhCN = ""
	this.Country.Names.ZhCN = ""
	if len(this.Subdivisions) > 0 {
		this.Subdivisions[0].Names.ZhCN = ""
	}
}

//Geoip2 地理位置解析结构体
type Geoip2 struct {
	mmdb       *maxminddb.Reader
	resultPool sync.Pool
}

//打开文件获取mmdb句柄
func NewGeoip(geoipMmdbByte []byte) (geoip *Geoip2, err error) {
	mmdb, err := maxminddb.FromBytes(geoipMmdbByte)
	if err != nil {
		return
	}
	return &Geoip2{mmdb: mmdb}, nil
}

func (this *Geoip2) Close() (err error) {
	return this.mmdb.Close()
}

func (this *Geoip2) Get() *Res {
	res := this.resultPool.Get()
	if res == nil {
		return &Res{}
	}
	return res.(*Res)
}

//通过ip获取地区
func (this *Geoip2) GetAreaFromIP(rawIP string) (province, city string, err error) {

	res := this.Get()
	res.reset()

	defer this.resultPool.Put(res)

	ip := net.ParseIP(rawIP)
	if ip == nil {
		logs.Logger.Error("net.ParseIP", zap.String("can't parse ip", rawIP))
		return "未知", "未知", nil
	}

	err = this.mmdb.Lookup(ip, &res)
	if err != nil {
		logs.Logger.Error("this.mmdb.Lookup", zap.String("err", err.Error()))
		return "未知", "未知", nil
	}

	if len(res.Subdivisions) > 0 {
		province = res.Subdivisions[0].Names.ZhCN
	}

	city = res.City.Names.ZhCN
	return province, city, nil
}
