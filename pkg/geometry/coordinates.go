// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

/**
** This file contains helper functions for coordinate data
**/
package geometry

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/tidwall/geodesic"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
)

const (
	GEODESIC_DIST_THRESHOLD = 100 // getting as close as this threshold (in meters) to a coordinate is sufficiently exact (111km ~= 1° longitude at equator)
	MAX_ITERATIONS          = 20
)

// FormatPolygonArray formats a given input polygon for usage in postGIS/SQL syntax
// Input is 2-dimensional array of points formatted: [[long1,lat1],[long2,lat2],...]
// Output is postGIS geometry syntax: (long1 lat1, long2 lat2, ...)
func FormatPolygonArray(polygon []model.SimplePoint) (string, error) {
	output := "("
	for i, point := range polygon {
		if i > 0 {
			// add separator before adding next point
			output += ","
		}
		output += fmt.Sprintf("%f %f", point.X, point.Y)
	}
	output += ")"
	return output, nil
}

// ParsePointArray parses a string representation of an array of float-points into an array of model.SimplePoint
func ParsePointArray(arrayString string) ([]model.SimplePoint, error) {
	polygon := []model.SimplePoint{}
	points := [][]float64{}
	err := json.Unmarshal([]byte(arrayString), &points)
	if err != nil {
		return nil, err
	}
	for _, p := range points {
		polygon = append(polygon, model.SimplePoint{X: p[0], Y: p[1]})
	}
	return polygon, nil
}

// CalcTranslation calculates the longitudinal translation and crossed bound of a polygon
// Input is 2-dimensional array of points formatted: [[long1,lat1],[long2,lat2],...]
// Max polygon size is limited to 180 height and 360 width. Scale polygon accordingly before these calculations.
// If the polygon is entirely within the X-axis bounds of -180 to +180, boundary and translationFactor are returned as 0.0
func CalcTranslation(polygon []model.SimplePoint) (float64, float64, error) {
	var left, right, top, bottom float64
	for i, point := range polygon {
		if i == 0 || left > point.X {
			left = point.X
		}
		if i == 0 || right < point.X {
			right = point.X
		}
		if i == 0 || top < point.Y {
			top = point.Y
		}
		if i == 0 || bottom > point.Y {
			bottom = point.Y
		}
	}
	if right-left > 360 || top-bottom > 180 {
		return 0, 0, fmt.Errorf("polygon dimensions out of bounds")
	}
	// since max width is 360, all points have the same translation (postGIS' wrapx handles cases where part of the bbox is within the bounds an thus should have a factor of 0)
	// and only one boundary can be crossed (-180 or +180)
	boundary := 0.0
	translationFactor := 0.0
	if right > 180 {
		boundary = 180
		translationFactor = -math.Floor((right + 180) / 360)
	}
	if left < -180 {
		boundary = -180.0
		translationFactor = -math.Floor((left + 180) / 360)
	}

	return boundary, translationFactor, nil
}

// WrapPolygon mimics wrapping a polygon around the +-180 degrees meridian
// returns a GeoJSONFeature consisting of the polygons parts cut and wrapped around the +-180 boundary
func WrapPolygonLon(polygon model.GeoJSONFeature) (*model.GeoJSONFeature, error) {
	// handle Point shapes
	if polygon.Geometry.Type == model.GEOJSON_GEOMETRY_POINT || polygon.Geometry.Type == model.GEOJSON_GEOMETRY_MULTIPOINT {
		// just translate the longitudes to be within -180 to +180
		for i, position := range polygon.Geometry.Coordinates {
			switch polygon.Geometry.Type {
			case model.GEOJSON_GEOMETRY_POINT:
				// coordinates are directly on the first level
				lon, ok := position.(float64)
				if !ok {
					return nil, fmt.Errorf("invalid GeoJSON coordinate type: %t", position)
				}
				polygon.Geometry.Coordinates[i] = TranslateLon(lon)
			case model.GEOJSON_GEOMETRY_MULTIPOINT:
				// coordinates are nested
				coords, ok := position.([]float64)
				if !ok {
					return nil, fmt.Errorf("invalid GeoJSON %s: %+v", polygon.Geometry.Type, polygon)
				}
				coords[0] = TranslateLon(coords[0])
				polygon.Geometry.Coordinates[i] = coords
			}
		}
		return &polygon, nil
	}
	// handle complex shapes that need cutting
	shapes, err := GetSimplePointShapes(&polygon.Geometry)
	if err != nil {
		return nil, err
	}
	partials := [][]model.SimplePoint{}
	for _, shape := range shapes {
		_, factor, err := CalcTranslation(shape)
		if err != nil {
			return nil, err
		}
		if factor == 0 {
			// polygon is within bounds, so leave it as-is
			partials = append(partials, shape)
			continue
		}
		part := cutGeometry(shape, polygon.Geometry.Type)
		partials = append(partials, part)
		translated := translatePolygonLon(shape, factor)
		part = cutGeometry(translated, polygon.Geometry.Type)
		partials = append(partials, part)
	}
	multi := ParseMultiPolygon(polygon.Geometry.Type, partials)
	bm, _ := json.Marshal(multi)
	orig, _ := json.Marshal(polygon)
	fmt.Printf("Original polygon:\n%+v\n", string(orig))
	fmt.Printf("Multipolygon:\n%+v\n", string(bm))
	return &multi, nil
}

// ParseMultiPolygon takes a list of polygons in the form of lists of model.SimplePoint and parses them as GeoJSON Multipolygon/MultiLineString
func ParseMultiPolygon(featureType model.GeoJSONGeometryType, parts [][]model.SimplePoint) model.GeoJSONFeature {
	if featureType == model.GEOJSON_GEOMETRY_MULTIPOLYGON || featureType == model.GEOJSON_GEOMETRY_POLYGON {
		// close shapes
		for i, part := range parts {
			if !IsClosed(part) {
				part = append(part, part[0])
				parts[i] = part
			}
		}
		polygons := []any{}
		for _, part := range parts {
			// polygons can contain multiple shapes
			shapes := [][][2]float64{}
			// a shape is a set of coordinate-tuples
			shape := [][2]float64{}
			for _, point := range part {
				shape = append(shape, [2]float64{point.X, point.Y})
			}
			shapes = append(shapes, shape)
			// add polygon to multipolygon
			polygons = append(polygons, shapes)
		}
		return model.GeoJSONFeature{
			Type: model.GEOJSONTYPE_FEATURE,
			Geometry: model.Geometry{
				Type:        model.GEOJSON_GEOMETRY_MULTIPOLYGON,
				Coordinates: polygons,
			},
		}
	}
	if featureType == model.GEOJSON_GEOMETRY_LINESTRING || featureType == model.GEOJSON_GEOMETRY_MULTILINESTRING {
		lines := []any{}
		for _, part := range parts {
			line := [][2]float64{}
			for _, point := range part {
				line = append(line, [2]float64{point.X, point.Y})
			}
			lines = append(lines, line)
		}
		return model.GeoJSONFeature{
			Type: model.GEOJSONTYPE_FEATURE,
			Geometry: model.Geometry{
				Type:        model.GEOJSON_GEOMETRY_MULTILINESTRING,
				Coordinates: lines,
			},
		}
	}
	return model.GeoJSONFeature{Type: "UNSUPPORTED TYPE"}
}

// ParsePolygon takes a polygon as a list of model.SimplePoints and parses it as a GeoJSON Polygon
func ParsePolygon(polygon []model.SimplePoint) model.GeoJSONFeature {
	coordinates := []any{}
	for _, p := range polygon {
		coordinates = append(coordinates, []float64{p.X, p.Y})
	}
	// close polygon if it not already is
	if !(polygon[0].X == polygon[len(polygon)-1].X && polygon[0].Y == polygon[len(polygon)-1].Y) {
		coordinates = append(coordinates, []float64{polygon[0].X, polygon[0].Y})
	}
	return model.GeoJSONFeature{
		Type: model.GEOJSONTYPE_FEATURE,
		Geometry: model.Geometry{
			Type:        model.GEOJSON_GEOMETRY_POLYGON,
			Coordinates: []any{coordinates},
		},
	}
}

// GetSimplePointShapes gather coordinates as collection of model.SimplePoint arrays - this can be a line, polygon or multi-feature
func GetSimplePointShapes(geom *model.Geometry) ([][]model.SimplePoint, error) {
	shapes := [][]model.SimplePoint{}
	if len(geom.Coordinates) == 0 {
		return shapes, nil
	}
	if geom.Type == model.GEOJSON_GEOMETRY_LINESTRING {
		shape := []model.SimplePoint{}
		// coordinates are nested in coordinates
		for _, point := range geom.Coordinates {
			coords, ok := point.([]float64)
			if !ok {
				return nil, fmt.Errorf("invalid GeoJSON %s: %+v", geom.Type, geom)
			}
			shape = append(shape, model.SimplePoint{X: coords[0], Y: coords[1]})
		}
		shapes = append(shapes, shape)
		return shapes, nil
	}
	if geom.Type == model.GEOJSON_GEOMETRY_MULTILINESTRING || geom.Type == model.GEOJSON_GEOMETRY_POLYGON {
		// coordinates are nested in shapes
		for _, s := range geom.Coordinates {
			shape := []model.SimplePoint{}
			s, ok := s.([]any)
			if !ok {
				return nil, fmt.Errorf("invalid GeoJSON %s: %+v", geom.Type, geom)
			}
			for _, point := range s {
				point, ok := point.([]float64)
				if !ok {
					return nil, fmt.Errorf("invalid GeoJSON %s: %+v", geom.Type, geom)
				}
				shape = append(shape, model.SimplePoint{X: point[0], Y: point[1]})
			}
			shapes = append(shapes, shape)
		}
		return shapes, nil
	}
	return nil, nil
}

// cutGeometry cuts the given geometry based on its type
func cutGeometry(polygon []model.SimplePoint, typ model.GeoJSONGeometryType) []model.SimplePoint {
	if typ == model.GEOJSON_GEOMETRY_LINESTRING || typ == model.GEOJSON_GEOMETRY_MULTILINESTRING {
		return cutLine(polygon)
	} else {
		return cutPolygon(polygon)
	}
}

// cutPolygon cuts a polygon at the vertical boundaris +-180 and returns the part of the polygon that is inside the boundary
// It traverses the vertices of the polygon and when it detects a vertex out of bounds it cuts the polygon and inserts the bounds-crossing-point as a new vertex
// TODO: what if cuttting at the boundary results in multiple disconnected parts? -> can only happen on concave polygons!
func cutPolygon(polygon []model.SimplePoint) []model.SimplePoint {
	// "open" the polygon as a duplicate point messes with the algorithm - we close the cut polygon again in the end
	wasClosed := false
	if IsClosed(polygon) {
		polygon = polygon[:len(polygon)-1]
		wasClosed = true
	}
	partial := []model.SimplePoint{}
	for i, p := range polygon {
		var lastPoint model.SimplePoint
		if i > 0 {
			lastPoint = polygon[i-1]
		} else {
			lastPoint = polygon[len(polygon)-1]
		}
		if outOfBounds(p.X) != outOfBounds(lastPoint.X) {
			// exactly one of the vertices is out of bounds - take the crossed boundary from that vertex
			var boundary float64
			if outOfBounds(p.X) {
				boundary = crossedBoundary(p.X)
			} else {
				boundary = crossedBoundary(lastPoint.X)
			}
			// add a crossing point
			crossingPoint := model.SimplePoint{X: boundary, Y: calcCrossingYCartesian(p, lastPoint, boundary)}
			partial = append(partial, crossingPoint)
		}
		if !outOfBounds(p.X) {
			// keep points that are in bounds
			partial = append(partial, p)
		}
	}
	// close polygon again
	if wasClosed && !IsClosed(partial) {
		partial = append(partial, partial[0])
	}
	return partial
}

// cutLine cuts a line at the vertical boundaris +-180 and returns the part of the line that is inside the boundary
// It traverses the vertices of the line and when it detects a vertex out of bounds it cuts the line and inserts the bounds-crossing-point as a new vertex
// TODO: what if cuttting at the boundary results in multiple disconnected parts? -> can only happen on concave polygons!
func cutLine(line []model.SimplePoint) []model.SimplePoint {
	partial := []model.SimplePoint{}
	for i, p := range line {
		var lastPoint model.SimplePoint
		if i > 0 {
			lastPoint = line[i-1]
		}
		// for lines, crossing can only happen for i >= 1
		if i > 0 && outOfBounds(p.X) != outOfBounds(lastPoint.X) {
			// exactly one of the vertices is out of bounds - take the crossed boundary from that vertex
			var boundary float64
			if outOfBounds(p.X) {
				boundary = crossedBoundary(p.X)
			} else {
				boundary = crossedBoundary(lastPoint.X)
			}
			// add a crossing point
			crossingPoint := model.SimplePoint{X: boundary, Y: calcCrossingYCartesian(p, lastPoint, boundary)}
			partial = append(partial, crossingPoint)
		}
		if !outOfBounds(p.X) {
			// keep points that are in bounds
			partial = append(partial, p)
		}
	}
	return partial
}

// outOfBounds returns whether an x (longitude) value lies out of the boundaries (+ or - 180)
func outOfBounds(x float64) bool {
	return x < -180 || x > 180
}

// crossedBoundary returns the boundary that has been crossed by a points X value: -180 or +180
// returns the X value if its inside the boundaries
func crossedBoundary(x float64) float64 {
	if x < -180 {
		return -180
	}
	if x > 180 {
		return 180
	}
	return x
}

// IsClosed returns whether a polygon ha a closed shape, meaning that its first and last point are the same
func IsClosed(polygon []model.SimplePoint) bool {
	if len(polygon) < 2 {
		// any polygon with 0 or 1 vertices cannot be closed
		return false
	}
	lastPos := len(polygon) - 1
	return polygon[0].X == polygon[lastPos].X && polygon[0].Y == polygon[lastPos].Y
}

// translatePolygonLon translates a polygon in longitudinal direction by a multiple of +/-360 degrees longitude given by the calculated factor
// So any vertices that are outside the bounds in the orignal polygon are now inside the bounds in the translated copy (and ones that were within the bounds are pushed out on the other side)
func translatePolygonLon(polygon []model.SimplePoint, factor float64) []model.SimplePoint {
	moved := make([]model.SimplePoint, 0, len(polygon))
	for _, point := range polygon {
		moved = append(moved, model.SimplePoint{X: point.X + (factor * 360), Y: point.Y})
	}
	return moved
}

// returns longitude coordinate translated to be within -180 to +180
func TranslateLon(lon float64) float64 {
	if lon <= 180 && lon >= -180 {
		return lon
	}
	if lon < 0 {
		for lon < -180 {
			lon += 360
		}
	} else {
		for lon > 180 {
			lon -= 360
		}
	}
	return lon
}

// calcCrossingYCartesian calculates the Y (Lat) coordinate of the point where the line between p1 and p2 crosses the x-boundary
func calcCrossingYCartesian(p1 model.SimplePoint, p2 model.SimplePoint, xBoundary float64) float64 {
	if p1.X == p2.X {
		// if x values are the same, m goes to infinity and every point on the line between p1 and p2 is a crossing point
		return p1.Y
	}
	m := (p1.Y - p2.Y) / (p1.X - p2.X)
	b := p1.Y - m*p1.X
	return m*xBoundary + b
}

// calcCrossingYGeodesic uses geodesic algorithms to approximate the crossing point of the geodesic line p1 - p2 and the x-boundary
// NOTE: This iterative approximation is still not significantly better than the carthesian approach in terms of fitting the orignal polygon line
func calcCrossingYGeodesic(p1 model.SimplePoint, p2 model.SimplePoint, xBoundary float64) float64 {
	var azimuth, crossY, crossX float64
	// calculate azimuth of geodesic line between p1 and p2
	geodesic.WGS84.Inverse(p1.Y, p1.X, p2.Y, p2.X, nil, &azimuth, nil)
	// get iteratively closer to the crossing point by "walking" from p1 in azimuth direction until we are reasonably close to long=xBoundary
	crossX = p1.X
	crossY = p1.Y
	var step float64
	for range MAX_ITERATIONS {
		// use delta X as step size to get closer to xBoundary
		geodesic.WGS84.Inverse(crossY, crossX, crossY, xBoundary, &step, nil, nil)
		step = math.Abs(step)
		if step < GEODESIC_DIST_THRESHOLD {
			break
		}
		// calculate coordinates of point from p1 in 'step' distance in azimuth direction
		geodesic.WGS84.Direct(crossY, crossX, azimuth, step, &crossY, &crossX, &azimuth)
	}
	return crossY
}
