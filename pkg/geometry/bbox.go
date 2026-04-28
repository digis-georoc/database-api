// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

/**
** This file contains helper functions for handling bounding boxes
**/
package geometry

import (
	"math"

	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
)

const (
	LONG_MIN = -180.0
	LONG_MAX = 180.0
	LAT_MIN  = -90.0
	LAT_MAX  = 90.0
)

// IsZoom0 returns true if the given bbox is big enough to fit the whole world view (zoom = 0); false if it is smaller.
func IsZoom0(bbox []model.SimplePoint) bool {
	// check if height of bbox is >= |LAT_MIN| + LAT_MAX and width of bbox is |LONG_MIN| + LONG_MAX
	return bbox[3].Y-bbox[0].Y >= math.Abs(LAT_MIN)+LAT_MAX && bbox[1].X-bbox[0].X >= math.Abs(LONG_MIN)+LONG_MAX
}

// ScaleBbox takes an array of coordinates for a bounding box and scales it around the center
// Scales latitudes up to LAT_MIN/LAT_MAX only
// Input is 2-dimensional array of points formatted: [[long1,lat1],[long2,lat2],...]
func ScaleBBox(bbox []model.SimplePoint) []model.SimplePoint {
	// calc width and height of bbox
	width := bbox[1].X - bbox[0].X
	height := bbox[3].Y - bbox[0].Y
	scaleLong := width / 2
	scaleLat := height / 2
	// add half bbox on each side
	// SW
	bbox[0].X = bbox[0].X - scaleLong
	bbox[0].Y = math.Max(LAT_MIN, bbox[0].Y-scaleLat)
	// SE
	bbox[1].X = bbox[1].X + scaleLong
	bbox[1].Y = math.Max(LAT_MIN, bbox[1].Y-scaleLat)
	// NE
	bbox[2].X = bbox[2].X + scaleLong
	bbox[2].Y = math.Min(LAT_MAX, bbox[2].Y+scaleLat)
	// NW
	bbox[3].X = bbox[3].X - scaleLong
	bbox[3].Y = math.Min(LAT_MAX, bbox[3].Y+scaleLat)
	return bbox
}

// TruncateBBox truncates a given bbox to at most 360 width and 180 height by reducing the size equally from both sides
func TruncateBBox(bbox []model.SimplePoint) []model.SimplePoint {
	width := bbox[1].X - bbox[0].X
	height := bbox[3].Y - bbox[0].Y
	if width <= LONG_MAX*2 && height <= LAT_MAX*2 {
		// no need to truncate
		return bbox
	}
	// calc middle longitude and latitude
	middleLong := (bbox[1].X + bbox[0].X) / 2
	middleLat := (bbox[3].Y + bbox[0].Y) / 2
	// truncate bbox by keeping only middle +- LONG/LAT_MAX
	// SW
	bbox[0].X = math.Max(bbox[0].X, middleLong-LONG_MAX)
	bbox[0].Y = math.Max(bbox[0].Y, middleLat-LAT_MAX)
	// SE
	bbox[1].X = math.Min(bbox[1].X, middleLong+LONG_MAX)
	bbox[1].Y = math.Max(bbox[1].Y, middleLat-LAT_MAX)
	// NE
	bbox[2].X = math.Min(bbox[2].X, middleLong+LONG_MAX)
	bbox[2].Y = math.Min(bbox[2].Y, middleLat+LAT_MAX)
	// NW
	bbox[3].X = math.Max(bbox[3].X, middleLong-LONG_MAX)
	bbox[3].Y = math.Min(bbox[3].Y, middleLat+LAT_MAX)
	return bbox
}
