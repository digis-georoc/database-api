package model

type Specimen struct {
	SamplingFeatureID int    `json:"samplingFeatureID"`
	SpecimenTypeCV    string `json:"specimenTypeCV"`
	SpecimenMediumCV  string `json:"specimenMediumCV"`
	IsFieldSpecimen   bool   `json:"isFieldSpecimen"`
}

type SpecimenResponse struct {
	NumItems int        `json:"numItems"`
	Data     []Specimen `json:"data"`
}

type Sample struct {
	SamplingFeatureID          int     `json:"samplingFeatureID"`
	SamplingFeatureUUID        string  `json:"samplingFeatureUUID"`
	SamplingFeatureName        string  `json:"samplingFeatureName"`
	SamplingFeatureDescription string  `json:"samplingFeatureDescription"`
	SamplingFeatureGeotypeCV   string  `json:"samplingFeatureGeoTypeCV"`
	FeatureGeometryWKT         string  `json:"featureGeometryWKT"`
	Elevation_m                float64 `json:"elevation_m"`
	ElevationDatumCV           string  `json:"elevationDatumCV"`
	ElevationPrecision         float64 `json:"elevationPrecision"`
	ElevationPrecisionComment  string  `json:"elevationPrecisionComment"`
}

type SampleResponse struct {
	NumItems int      `json:"numItems"`
	Data     []Sample `json:"data"`
}

type SampleByFilters struct {
	SampleID  int     `json:"sampleID"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SampleByFilterResponse struct {
	NumItems int               `json:"numItems"`
	Data     []SampleByFilters `json:"data"`
}

type FilteredSample struct {
	ValuesString string `json:"valuesString"`
	NumSamples   int    `json:"numSamples"`
}

type ClusteredSample struct {
	ClusterID    int      `json:"clusterID"`
	Centroid     Geometry `json:"centroid"`
	ConvexHull   Geometry `json:"convexHull"`
	PointStrings []string `json:"pointsWithIds"`
	Samples      []int64  `json:"samples"`
}

type ClusterResponse struct {
	Clusters []GeoJSONCluster `json:"clusters"`
	Bbox     GeoJSONFeature   `json:"bbox"`
	Points   []GeoJSONFeature `json:"points"`
}

type SamplingTechnique struct {
	Name string
}

type SamplingTechniqueResponse struct {
	NumItems int                 `json:"numItems"`
	Data     []SamplingTechnique `json:"data"`
}

type Material struct {
	Name string
}

type MaterialResponse struct {
	NumItems int        `json:"numItems"`
	Data     []Material `json:"data"`
}

type InclusionType struct {
	Name string
}

type InclusionTypeResponse struct {
	NumItems int             `json:"numItems"`
	Data     []InclusionType `json:"data"`
}

type GeoAge struct {
	Name string
}

type GeoAgeResponse struct {
	NumItems int      `json:"numItems"`
	Data     []GeoAge `json:"data"`
}

type GeoAgePrefix struct {
	Name string
}

type GeoAgePrefixResponse struct {
	NumItems int            `json:"numItems"`
	Data     []GeoAgePrefix `json:"data"`
}

type Organization struct {
	Name string
}

type OrganizationResponse struct {
	NumItems int            `json:"numItems"`
	Data     []Organization `json:"data"`
}
