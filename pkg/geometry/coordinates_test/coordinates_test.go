// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package coordinates_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gitlab.gwdg.de/fe/digis/database-api/pkg/geometry"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
)

func TestTranslateLon(t *testing.T) {
	tests := map[float64]float64{
		-180:    -180,
		180:     180,
		0:       0,
		-52.37:  -52.37,
		169.005: 169.005,
		-181:    179,
		181:     -179,
		-365:    -5,
		365:     5,
	}
	for i, exp := range tests {
		o := geometry.TranslateLon(i)
		if o != exp {
			t.Fatalf("Input: %f | Output: %f | Expected: %f", i, o, exp)
		}
	}
}

func TestWrapPoint(t *testing.T) {
	point := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type:        model.GEOJSON_GEOMETRY_POINT,
			Coordinates: []any{-185.5, 89.0},
		},
	}
	expected := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type:        model.GEOJSON_GEOMETRY_POINT,
			Coordinates: []any{174.5, 89.0},
		},
	}
	wrapped, err := geometry.WrapPolygonLon(point)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", expected) != fmt.Sprintf("%+v", *wrapped) {
		t.Fatalf("Expected:\n%+v\n but got:\n%+v", expected, *wrapped)
	}
}

func TestWrapMultiPoint(t *testing.T) {
	point := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_MULTIPOINT,
			Coordinates: []any{
				[]any{
					-177.3,
					41.1,
				},
				[]any{
					197.0,
					-4.3,
				},
				[]any{
					-700.0,
					38.0,
				},
			},
		},
	}
	expected := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_MULTIPOINT,
			Coordinates: []any{
				[]any{
					-177.3,
					41.1,
				},
				[]any{
					-163,
					-4.3,
				},
				[]any{
					20,
					38,
				},
			},
		},
	}
	wrapped, err := geometry.WrapPolygonLon(point)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", expected) != fmt.Sprintf("%+v", *wrapped) {
		t.Fatalf("Expected:\n%+v\n but got:\n%+v", expected, *wrapped)
	}
}

func TestWrapLine(t *testing.T) {
	line := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_LINESTRING,
			Coordinates: []any{
				[]any{
					-197.0,
					20.0,
				},
				[]any{
					-170.,
					-20.,
				},
			},
		},
	}
	expected := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_MULTILINESTRING,
			Coordinates: []any{
				[]any{ // first line
					[]any{
						-180,
						-5.185185185185219,
					},
					[]any{
						-170,
						-20,
					},
				},
				[]any{ // second line
					[]any{
						163,
						20,
					},
					[]any{
						180,
						-5.185185185185162,
					},
				},
			},
		},
	}
	wrapped, err := geometry.WrapPolygonLon(line)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", expected) != fmt.Sprintf("%+v", *wrapped) {
		t.Fatalf("Expected:\n%+v\n but got:\n%+v", expected, *wrapped)
	}
}

func TestWrapPolygonRect(t *testing.T) {
	polygon := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_POLYGON,
			Coordinates: []any{ // outer array
				[]any{ // polygon-parts
					[]any{ // coordinate level
						160.,
						10.,
					},
					[]any{
						200.,
						10.,
					},
					[]any{
						200.,
						0.,
					},
					[]any{
						160.,
						0.,
					},
				},
			},
		},
	}
	expected := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_MULTIPOLYGON,
			Coordinates: []any{ // outer array
				[]any{ // polygon
					[]any{ // polygon-parts
						[]any{ // coordinate level
							160,
							10,
						},
						[]any{
							180,
							10,
						},
						[]any{
							180,
							0,
						},
						[]any{
							160,
							0,
						},
						[]any{ // closing vertex
							160,
							10,
						},
					},
				},
				[]any{ // polygon2
					[]any{ // polygon-parts
						[]any{ // coordinate level
							-180,
							10,
						},
						[]any{
							-160,
							10,
						},
						[]any{
							-160,
							0,
						},
						[]any{
							-180,
							0,
						},
						[]any{ // closing vertex
							-180,
							10,
						},
					},
				},
			},
		},
	}
	wrapped, err := geometry.WrapPolygonLon(polygon)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", expected) != fmt.Sprintf("%+v", *wrapped) {
		t.Fatalf("Expected:\n%+v\n but got:\n%+v", expected, *wrapped)
	}
}

func TestWrapPolygonConcave(t *testing.T) {
	polygon := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_POLYGON,
			Coordinates: []any{ // outer array
				[]any{ // polygon-parts
					[]any{ // coordinate level
						160.,
						10.,
					},
					[]any{
						200.,
						5.,
					},
					[]any{
						160.,
						0.,
					},
					[]any{
						200.,
						-5.,
					},
					[]any{
						160.,
						-10.,
					},
					[]any{
						160.,
						10.,
					},
				},
			},
		},
	}
	expected := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_MULTIPOLYGON,
			Coordinates: []any{ // outer array
				[]any{ // polygon
					[]any{ // polygon-parts
						[]any{ // coordinate level
							160,
							10,
						},
						[]any{
							180,
							7.5,
						},
						[]any{
							180,
							2.5,
						},
						[]any{
							160,
							0,
						},
						[]any{
							180,
							-2.5,
						},
						[]any{
							180,
							-7.5,
						},
						[]any{
							160,
							-10,
						},
						[]any{
							160,
							10,
						},
					},
				},
				[]any{ // polygon2
					[]any{ // polygon-parts
						[]any{ // coordinate level
							-180,
							7.5,
						},
						[]any{
							-160,
							5,
						},
						[]any{
							-180,
							2.5,
						},
						[]any{
							-180,
							-2.5,
						},
						[]any{
							-160,
							-5,
						},
						[]any{
							-180,
							-7.5,
						},
						[]any{
							-180,
							7.5,
						},
					},
				},
			},
		},
	}
	wrapped, err := geometry.WrapPolygonLon(polygon)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", expected) != fmt.Sprintf("%+v", *wrapped) {
		t.Fatalf("Expected:\n%+v\n but got:\n%+v", expected, *wrapped)
	}
}

func TestGeoJSONPolygons(t *testing.T) {
	path := "./polygons"
	files, err := os.ReadDir(path)
	if err != nil {
		t.Fatalf("Can not read directory: %s", err.Error())
	}
	for _, f := range files {
		b, err := os.ReadFile(path + "/" + f.Name())
		if err != nil {
			t.Fatalf("Can not read file: %s", err.Error())
		}
		col := model.GeoJSONFeatureCollection{}
		err = json.Unmarshal(b, &col)
		if err != nil {
			t.Fatalf("Can not unmarshal file as GeoJSONCollection: %s", err.Error())
		}
		for _, feat := range col.Features {
			wrapped, err := geometry.WrapPolygonLon(feat)
			if err != nil {
				t.Fatalf("Can not wrap polygon: %s", err.Error())
			}
			if len(wrapped.Geometry.Coordinates) == 0 {
				t.Fatalf("Empty result")
			}
			// check if all coordinates are in bounds +-180 lon and +-90 lat
			shapes, err := geometry.GetSimplePointShapes(&wrapped.Geometry)
			if err != nil {
				t.Fatalf("Can not get geometry as simple points: %s", err.Error())
			}
			for _, s := range shapes {
				for _, p := range s {
					if p.X < -180 || p.X > 180 || p.Y < -90 || p.Y > 90 {
						t.Fatalf("Point out of bounds: %+v", p)
					}
				}
			}
		}
	}
}
