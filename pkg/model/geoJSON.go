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
	Latitude              *float64 `json:"latitude"`
	Longitude             *float64 `json:"longitude"`
	LocationID            *int     `json:"locationID"`
	NumSamplingFeatureIDs *int     `json:"numSamplingFeatureIDs"`
	SamplingFeatureIDs    []*int   `json:"samplingFeatureIDs"`
	Setting               *string  `json:"setting"`
	Loc1                  *string  `json:"loc1"`
	Loc2                  *string  `json:"loc2"`
	Loc3                  *string  `json:"loc3"`
	LandOrSea             *string  `json:"landOrSea"`
}
