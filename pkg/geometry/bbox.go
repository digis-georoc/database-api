// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

/**
** This file contains helper functions for handling bounding boxes
**/
package geometry

import "math"

const (
	LONG_MIN = -180.0
	LONG_MAX = 180.0
	LAT_MIN  = -90.0
	LAT_MAX  = 90.0
)

// IsZoom0 returns true if the given bbox is big enough to fit the whole world view (zoom = 0); false if it is smaller.
func IsZoom0(bbox [][]float64) bool {
	// check if height of bbox is >= |LAT_MIN| + LAT_MAX and width of bbox is |LONG_MIN| + LONG_MAX
	return bbox[3][1]-bbox[0][1] >= math.Abs(LAT_MIN)+LAT_MAX && bbox[1][0]-bbox[0][0] >= math.Abs(LONG_MIN)+LONG_MAX
}

// ScaleBbox takes an array of coordinates for a bounding box and scales it around the center
// Scales latitudes up to LAT_MIN/LAT_MAX only
// Input is 2-dimensional array of points formatted: [[long1,lat1],[long2,lat2],...]
func ScaleBBox(bbox [][]float64) [][]float64 {
	// calc width and height of bbox
	width := bbox[1][0] - bbox[0][0]
	height := bbox[3][1] - bbox[0][1]
	scaleLong := width / 2
	scaleLat := height / 2
	// add half bbox on each side
	// SW
	bbox[0][0] = bbox[0][0] - scaleLong
	bbox[0][1] = math.Max(LAT_MIN, bbox[0][1]-scaleLat)
	// SE
	bbox[1][0] = bbox[1][0] + scaleLong
	bbox[1][1] = math.Max(LAT_MIN, bbox[1][1]-scaleLat)
	// NE
	bbox[2][0] = bbox[2][0] + scaleLong
	bbox[2][1] = math.Min(LAT_MAX, bbox[2][1]+scaleLat)
	// NW
	bbox[3][0] = bbox[3][0] - scaleLong
	bbox[3][1] = math.Min(LAT_MAX, bbox[3][1]+scaleLat)
	return bbox
}

// TruncateBBox truncates a given bbox to at most 360 width and 180 height by reducing the size equally from both sides
func TruncateBBox(bbox [][]float64) [][]float64 {
	width := bbox[1][0] - bbox[0][0]
	height := bbox[3][1] - bbox[0][1]
	if width <= LONG_MAX*2 && height <= LAT_MAX*2 {
		// no need to truncate
		return bbox
	}
	// calc middle longitude and latitude
	middleLong := (bbox[1][0] + bbox[0][0]) / 2
	middleLat := (bbox[3][1] + bbox[0][1]) / 2
	// truncate bbox by keeping only middle +- LONG/LAT_MAX
	// SW
	bbox[0][0] = math.Max(bbox[0][0], middleLong-LONG_MAX)
	bbox[0][1] = math.Max(bbox[0][1], middleLat-LAT_MAX)
	// SE
	bbox[1][0] = math.Min(bbox[1][0], middleLong+LONG_MAX)
	bbox[1][1] = math.Max(bbox[1][1], middleLat-LAT_MAX)
	// NE
	bbox[2][0] = math.Min(bbox[2][0], middleLong+LONG_MAX)
	bbox[2][1] = math.Min(bbox[2][1], middleLat+LAT_MAX)
	// NW
	bbox[3][0] = math.Max(bbox[3][0], middleLong-LONG_MAX)
	bbox[3][1] = math.Min(bbox[3][1], middleLat+LAT_MAX)
	return bbox
}
