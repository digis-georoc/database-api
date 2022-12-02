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
	Type           GeoJSONFeatureType
	Features       []GeoJSONFeature
	NumberMatched  int
	NumberReturned int
}

// GeoJSON Geometry
type Geometry struct {
	Type        GeoJSONFeatureType
	Coordinates []interface{}
}

// GeoJSON Feature
type GeoJSONFeature struct {
	Type       GeoJSONGeometryType
	ID         string
	Geometry   Geometry
	Properties map[string]interface{}
}
