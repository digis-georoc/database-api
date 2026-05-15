// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package coordinates_test

import (
	"fmt"
	"reflect"
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
	if !reflect.DeepEqual(*wrapped, expected) {
		t.Fatalf("Expected:\n%+v\n but got:\n%+v", expected, *wrapped)
	}
}

func TestWrapMultiPoint(t *testing.T) {
	point := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_MULTIPOINT,
			Coordinates: []any{
				[]float64{
					-177.3,
					41.1,
				},
				[]float64{
					197,
					-4.3,
				},
				[]float64{
					-700,
					38,
				},
			},
		},
	}
	expected := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_MULTIPOINT,
			Coordinates: []any{
				[]float64{
					-177.3,
					41.1,
				},
				[]float64{
					-163,
					-4.3,
				},
				[]float64{
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
	if !reflect.DeepEqual(*wrapped, expected) {
		t.Fatalf("Expected:\n%+v\n but got:\n%+v", expected, *wrapped)
	}
}

func TestWrapLine(t *testing.T) {
	point := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_LINESTRING,
			Coordinates: []any{
				[]float64{
					-197,
					20,
				},
				[]float64{
					-170,
					-20,
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
					[]float64{
						-180,
						-5.185185185185219,
					},
					[]float64{
						-170,
						-20,
					},
				},
				[]any{ // second line
					[]float64{
						163,
						202,
					},
					[]float64{
						180,
						-5.185185185185162,
					},
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

func TestWrapPolygon(t *testing.T) {
	point := model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type: model.GEOJSON_GEOMETRY_POLYGON,
			Coordinates: []any{ // outer array
				[]any{ // polygon-parts
					[]float64{ // coordinate level
						160,
						10,
					},
					[]float64{
						200,
						10,
					},
					[]float64{
						200,
						0,
					},
					[]float64{
						160,
						0,
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
						[]float64{ // coordinate level
							160,
							10,
						},
						[]float64{
							180,
							10,
						},
						[]float64{
							180,
							0,
						},
						[]float64{
							160,
							0,
						},
						[]float64{ // closing vertex
							160,
							10,
						},
					},
				},
				[]any{ // polygon2
					[]any{ // polygon-parts
						[]float64{ // coordinate level
							-180,
							10,
						},
						[]float64{
							-160,
							10,
						},
						[]float64{
							-160,
							0,
						},
						[]float64{
							-180,
							0,
						},
						[]float64{ // closing vertex
							-180,
							10,
						},
					},
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
