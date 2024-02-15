// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

type SpecimenType struct {
	SpecimenTypeCV string `json:"specimentypecv"`
}

type SpecimenTypeResponse struct {
	NumItems int            `json:"numItems"`
	Data     []SpecimenType `json:"data"`
}

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
	SamplingFeatureID          int      `json:"samplingFeatureID"`
	SamplingFeatureUUID        *string  `json:"samplingFeatureUUID"`
	SamplingFeatureName        *string  `json:"samplingFeatureName"`
	SamplingFeatureDescription *string  `json:"samplingFeatureDescription"`
	SamplingFeatureGeotypeCV   *string  `json:"samplingFeatureGeoTypeCV"`
	FeatureGeometryWKT         *string  `json:"featureGeometryWKT"`
	ElevationM                 *float64 `json:"elevationM"`
	ElevationDatumCV           *string  `json:"elevationDatumCV"`
	ElevationPrecision         *float64 `json:"elevationPrecision"`
	ElevationPrecisionComment  *string  `json:"elevationPrecisionComment"`
}

type SampleResponse struct {
	NumItems int      `json:"numItems"`
	Data     []Sample `json:"data"`
}

type SampleByFilters struct {
	SampleID          int     `json:"sampleID"`
	SampleName        string  `json:"sampleName"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Mineral           *string `json:"mineral"`
	RockClass         *string `json:"rockClass"`
	InclusionType     *string `json:"inclusionType"`
	GeologicalSetting *string `json:"geologicalSetting"`
	GeologicalAge     *string `json:"geologicalAge"`
	TotalCount        int     `json:"totalCount"`
}

type SampleByFiltersData struct {
	SampleID          int     `json:"sampleID"`
	SampleName        string  `json:"sampleName"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Mineral           *string `json:"mineral"`
	RockClass         *string `json:"rockClass"`
	InclusionType     *string `json:"inclusionType"`
	GeologicalSetting *string `json:"geologicalSetting"`
	GeologicalAge     *string `json:"geologicalAge"`
}

type SampleByFilterResponse struct {
	NumItems   int                   `json:"numItems"`
	TotalCount int                   `json:"totalCount"`
	Data       []SampleByFiltersData `json:"data"`
}

type ClusteredSample struct {
	ClusterID        int      `json:"clusterID"`
	CentroidString   string   `json:"centroid"`
	ConvexHullString string   `json:"convexHull"`
	PointStrings     []string `json:"points"`
	Samples          []int64  `json:"samples"`
}

type ClusterResponse struct {
	Clusters []GeoJSONCluster `json:"clusters"`
	Points   []GeoJSONFeature `json:"points"`
	Bbox     GeoJSONFeature   `json:"bbox"`
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
