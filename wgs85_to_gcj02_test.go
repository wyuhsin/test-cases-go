package tests

import (
	"math"
	"testing"
)

const (
	PI = 3.1415926535897932384626
	ee = 0.00669342162296594323
)

func TestLocation(t *testing.T) {
	lat, lng := wgs84ToGcj02(36.06623, 120.38299)
	t.Logf("lat: %f, lng: %f\n", lat, lng)
}

func wgs84ToGcj02(lng float64, lat float64) (float64, float64) {
	a := 6378245.0
	dlat := transformlat(lng-105.0, lat-35.0)
	dlng := transformlng(lng-105.0, lat-35.0)
	radlat := (lat / 180.0) * PI
	magic := math.Sin(radlat)
	magic = 1 - ee*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * PI)
	dlng = (dlng * 180.0) / (a / sqrtmagic * math.Cos(radlat) * PI)
	mglat := lat + dlat
	mglng := lng + dlng
	return mglng, mglat
}

func transformlat(lng float64, lat float64) float64 {
	ret := -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*PI) + 40.0*math.Sin((lat/3.0)*PI)) * 2.0 / 3.0
	ret += (160.0*math.Sin((lat/12.0)*PI) + 320*math.Sin((lat*PI)/30.0)) * 2.0 / 3.0
	return ret
}

func transformlng(lng float64, lat float64) float64 {
	var ret = 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += ((20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0) / 3.0
	ret += ((20.0*math.Sin(lng*PI) + 40.0*math.Sin((lng/3.0)*PI)) * 2.0) / 3.0
	ret += ((150.0*math.Sin((lng/12.0)*PI) + 300.0*math.Sin((lng/30.0)*PI)) * 2.0) / 3.0
	return ret
}
