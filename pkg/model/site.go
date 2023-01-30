package model

type Site struct {
	SamplingFeatureID        int
	SiteTypeCV               string
	Latitude                 float64
	Longitude                float64
	SpatialReferenceID       int
	LocationPrecision        float64
	LocationPrecisionComment string
	SiteDescription          string
	Setting                  string
}

type LandOrSea struct {
	Name string
}
