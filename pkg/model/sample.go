package model

type Sample struct {
	SamplingFeatureID          int
	SamplingFeatureUUID        int
	SamplingFeatureTypeCV      string
	SamplingFeatureCode        string
	SamplingFeatureName        string
	SamplingFeatureDescription string
	ElevationPrecision         float64
	ElevationPrecisionComment  string
}

type Specimen struct {
	SpecimenType string
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
