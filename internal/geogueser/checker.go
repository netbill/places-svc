package geogueser

import (
	"embed"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

//go:embed ukraine.geojson
var fs embed.FS

type Checker struct {
	mp orb.MultiPolygon
}

func New() (Checker, error) {
	b, err := fs.ReadFile("ukraine.geojson")
	if err != nil {
		return Checker{}, err
	}

	fc, err := geojson.UnmarshalFeatureCollection(b)
	if err != nil {
		return Checker{}, err
	}

	var mp orb.MultiPolygon
	for _, f := range fc.Features {
		if f == nil || f.Geometry == nil {
			continue
		}

		var g orb.Geometry
		switch v := any(f.Geometry).(type) {
		case orb.Geometry:
			g = v
		case *geojson.Geometry:
			g = v.Geometry()
		default:
			continue
		}

		switch v := g.(type) {
		case orb.MultiPolygon:
			mp = append(mp, v...)
		case orb.Polygon:
			mp = append(mp, v)
		}
	}

	if len(mp) == 0 {
		return Checker{}, fmt.Errorf("no polygons found in geojson")
	}

	return Checker{mp: mp}, nil
}

func (c Checker) ContainsLatLng(lat, lng float64) bool {
	return planar.MultiPolygonContains(c.mp, orb.Point{lng, lat})
}
