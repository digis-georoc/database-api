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
)

// FormatPolygonArray formats a given input polygon for usage in postGIS/SQL syntax
// Input is 2-dimensional array of points formatted: [[long1,lat1],[long2,lat2],...]
// Output is postGIS geometry syntax: (long1 lat1, long2 lat2, ...)
func FormatPolygonArray(polygon [][]float64) (string, error) {
	output := "("
	for i, point := range polygon {
		if i > 0 {
			// add separator before adding next point
			output += ","
		}
		for _, coordinate := range point {
			output += fmt.Sprintf(" %f", coordinate)
		}
	}
	output += ")"
	return output, nil
}

// ParsePointArray parses a string representation of an array of float-points into a 2-dimensional array
func ParsePointArray(arrayString string) ([][]float64, error) {
	polygon := [][]float64{}
	err := json.Unmarshal([]byte(arrayString), &polygon)
	if err != nil {
		return nil, err
	}
	return polygon, nil
}

// CalcTranslation calculates the longitudinal translation and crossed bound of a polygon
// Input is 2-dimensional array of points formatted: [[long1,lat1],[long2,lat2],...]
// Max polygon size is limited to 180 height and 360 width. Scale polygon accordingly before these calculations.
func CalcTranslation(polygon [][]float64) (float64, float64, error) {
	var left, right, top, bottom float64
	for i, point := range polygon {
		if i == 0 || left > point[0] {
			left = point[0]
		}
		if i == 0 || right < point[0] {
			right = point[0]
		}
		if i == 0 || top < point[1] {
			top = point[1]
		}
		if i == 0 || bottom > point[1] {
			bottom = point[1]
		}
	}
	if right-left > 360 || top-bottom > 180 {
		return 0, 0, fmt.Errorf("Polygon dimensions out of bounds")
	}
	// since max width is 360, all points have the same translation (postGIS' wrapx handles cases where part of the bbox is within the bounds an thus should have a factor of 0)
	// and only one boundary can be crossed (-180 or +180)
	boundary := 180.0
	translationFactor := -math.Floor((right + 180) / 360)
	if left < -180 {
		boundary = -180.0
		translationFactor = -math.Floor((left + 180) / 360)
	}

	return boundary, translationFactor, nil
}
