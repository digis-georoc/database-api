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

// GeoJSON Geometry
type Geometry struct {
	Type        GeoJSONFeatureType `json:"type"`
	Coordinates []interface{}      `json:"coordinates"`
}

// GeoJSON Feature
type GeoJSONFeature struct {
	Type       GeoJSONGeometryType    `json:"type"`
	ID         string                 `json:"id"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}
