// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

type GeologicalSetting struct {
	Setting string `json:"setting"`
}

type GeologicalSettingResponse struct {
	NumItems int                 `json:"numItems"`
	Data     []GeologicalSetting `json:"data"`
}

type Site struct {
	SamplingFeatureID        int      `json:"samplingFeatureID"`
	SiteTypeCV               *string  `json:"siteTypeCV"`
	Latitude                 *float64 `json:"latitude"`
	Longitude                *float64 `json:"longitude"`
	SpatialReferenceID       *int     `json:"spatialReferenceID"`
	LocationPrecision        *float64 `json:"locationPrecision"`
	LocationPrecisionComment *string  `json:"locationPrecisionComment"`
	SiteDescription          *string  `json:"siteDescription"`
	Setting                  *string  `json:"setting"`
}

type SiteResponse struct {
	NumItems int    `json:"numItems"`
	Data     []Site `json:"data"`
}
