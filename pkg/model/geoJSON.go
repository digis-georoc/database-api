// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

// Models for GeoJSON representation. See https://geojson.org/

type GeoJSONFeatureType = string
type GeoJSONGeometryType = string

const (
	GEOJSONTYPE_FEATURECOLLECTION    GeoJSONFeatureType  = "FeatureCollection"
	GEOJSONTYPE_FEATURE              GeoJSONFeatureType  = "Feature"
	GEOJSON_GEOMETRY_POINT           GeoJSONGeometryType = "Point"
	GEOJSON_GEOMETRY_LINESTRING      GeoJSONGeometryType = "LineString"
	GEOJSON_GEOMETRY_POLYGON         GeoJSONGeometryType = "Polygon"
	GEOJSON_GEOMETRY_MULTIPOINT      GeoJSONGeometryType = "MultiPoint"
	GEOJSON_GEOMETRY_MULTILINESTRING GeoJSONGeometryType = "MultiLineString"
	GEOJSON_GEOMETRY_MULTIPOLYGON    GeoJSONGeometryType = "MultiPolygon"
)

// GeoJSON FeatureCollection
type GeoJSONFeatureCollection struct {
	Type           GeoJSONFeatureType `json:"type"`
	Features       []GeoJSONFeature   `json:"features"`
	NumberMatched  int                `json:"numberMatched"`
	NumberReturned int                `json:"numberReturned"`
}

// GeoJSON Feature
type GeoJSONFeature struct {
	Type       GeoJSONFeatureType     `json:"type"`
	ID         string                 `json:"id"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// GeoJSON Geometry
type Geometry struct {
	Type        GeoJSONGeometryType `json:"type"`
	Coordinates []interface{}       `json:"coordinates"`
}

// GeoJSON Cluster of map locations
type GeoJSONCluster struct {
	ClusterID  int            `json:"clusterID"`
	Centroid   GeoJSONFeature `json:"centroid"`
	ConvexHull GeoJSONFeature `json:"convexHull"`
}

// GeoJSONSite
// note: pointer types allow for NULL values in sql
type GeoJSONSite struct {
	// nullable
	Latitude *float64 `json:"latitude"`
	// nullable
	Longitude *float64 `json:"longitude"`
	// nullable
	LocationID *int `json:"locationID"`
	// nullable
	NumSamplingFeatureIDs *int `json:"numSamplingFeatureIDs"`
	// nullable
	SamplingFeatureIDs []*int `json:"samplingFeatureIDs"`
	// nullable
	Setting *string `json:"setting"`
	// nullable
	Loc1 *string `json:"loc1"`
	// nullable
	Loc2 *string `json:"loc2"`
	// nullable
	Loc3 *string `json:"loc3"`
	// nullable
	LandOrSea *string `json:"landOrSea"`
}
