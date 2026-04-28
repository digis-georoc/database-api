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

	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
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
func WrapPolygonLon(polygon []model.SimplePoint) ([]model.SimplePoint, []model.SimplePoint, error) {
	polygon = TruncateBBox(polygon)
	boundary, _, err := CalcTranslation(polygon)
	if err != nil {
		return nil, nil, err
	}
	partial1 := cutPolygon(polygon)
	translated := translatePolygonLon(polygon, boundary)
	partial2 := cutPolygon(translated)
	return partial1, partial2, nil
}

// cutPolygon cuts a polygon at the vertical boundaris +-180 and returns the part of the polygon that is inside the boundary
// It traverses the vertices of the polygon and when it detects a vertex out of bounds it cuts the polygon and inserts the bounds-crossing-point as a new vertex
// TODO: what if cuttting at the boundary results in multiple disconnected parts? -> can only happen on concave polygons!
func cutPolygon(polygon []model.SimplePoint) []model.SimplePoint {
	// "open" the polygon as a duplicate point messes with the algorithm - we close the cut polygon again in the end
	if isClosed(polygon) {
		polygon = polygon[:len(polygon)-1]
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
			crossingPoint := model.SimplePoint{X: boundary, Y: calcCrossingY(p, lastPoint, boundary)}
			partial = append(partial, crossingPoint)
		}
		if !outOfBounds(p.X) {
			// keep points that are in bounds
			partial = append(partial, p)
		}
	}
	// close polygon again
	if !isClosed(partial) {
		partial = append(partial, partial[0])
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

// isClosed returns whether a polygon ha a closed shape, meaning that its first and last point are the same
func isClosed(polygon []model.SimplePoint) bool {
	if len(polygon) < 2 {
		// any polygon with 0 or 1 vertices cannot be closed
		return false
	}
	lastPos := len(polygon) - 1
	return polygon[0].X == polygon[lastPos].X && polygon[0].Y == polygon[lastPos].Y
}

// translatePolygonLon translates a polygon in longitudinal direction to get its "real coordinates" in the bounds -180/+180
// it returns a translated copy of the polygon that is moved by +/-360 degrees longitude
// So any vertices that are outside the bounds in the orignal polygon are now inside the bounds in the translated copy
func translatePolygonLon(polygon []model.SimplePoint, boundary float64) []model.SimplePoint {
	moved := make([]model.SimplePoint, 0, len(polygon))
	// move polygon in negative boundary direction
	// crossed western map boundary -> move east 360
	// crossed eastern map boundary -> move west 360
	for _, point := range polygon {
		moved = append(moved, model.SimplePoint{X: point.X - 2*boundary, Y: point.Y})
	}
	return moved
}

// calcCrossingY calculates the Y (Lat) corrdinate of the point where the line between p1 and p2 crosses the x-boundary
func calcCrossingY(p1 model.SimplePoint, p2 model.SimplePoint, xBoundary float64) float64 {
	if p1.X == p2.X {
		// if x values are the same, the crossing Y can only
		return p1.Y
	}
	m := math.Abs(p1.Y-p2.Y) / math.Abs(p1.X-p2.X)
	b := p1.Y - m*p1.X
	return m*xBoundary + b
}
