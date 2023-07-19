package model

type Specimen struct {
	SamplingFeatureID int
	SpecimenTypeCV    string
	SpecimenMediumCV  string
	IsFieldSpecimen   bool
}

type Sample struct {
	SamplingFeatureID          int
	SamplingFeatureUUID        string
	SamplingFeatureName        string
	SamplingFeatureDescription string
	SamplingFeatureGeotypeCV   string
	FeatureGeometryWKT         string
	Elevation_m                float64
	ElevationDatumCV           string
	ElevationPrecision         float64
	ElevationPrecisionComment  string
}

type SampleByGeoSettingResponse struct {
	SamplingFeatureID int
	Specimen          int
	Latitude          float64 `json:"lat"`
	Longitude         float64 `json:"long"`
	Setting           string
	Location1         string `json:"loc1"`
	Location2         string `json:"loc2"`
	Location3         string `json:"loc3"`
	Texture           string
	RockType          string `json:"rock_type"`
	RockClass         string `json:"rock_class"`
	Mineral           string
	Material          string
	InclusionType     string   `json:"inclusion_type"`
	SamplingTechnique string   `json:"samp_technique"`
	SampleNames       []string `json:"sample_names"`
	LandOrSea         string   `json:"land_or_sea"`
	RimOrCore         string   `json:"rim_or_core"`
}

type SampleByFilterResponse struct {
	NumItems int               `json:"numItems"`
	Data     []SampleByFilters `json:"data"`
}

type SampleByFilters struct {
	SampleID  int     `json:"sampleID"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ClusteredSample struct {
	ClusterID  int      `json:"clusterid"`
	Centroid   Geometry `json:"centroid"`
	ConvexHull Geometry `json:"convexHull"`
	Points     Geometry `json:"points"`
	Samples    []int64  `json:"samples"`
}

type SamplingTechnique struct {
	Name string
}

type Material struct {
	Name string
}

type InclusionType struct {
	Name string
}

type GeoAge struct {
	Name string
}

type GeoAgePrefix struct {
	Name string
}

type Organization struct {
	Name string
}
