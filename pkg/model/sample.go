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
