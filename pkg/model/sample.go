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
	SamplingFeatureID int `json:"samplingFeatureID"`
	// nullable
	SamplingFeatureUUID *string `json:"samplingFeatureUUID"`
	// nullable
	SamplingFeatureName *string `json:"samplingFeatureName"`
	// nullable
	SamplingFeatureDescription *string `json:"samplingFeatureDescription"`
	// nullable
	SamplingFeatureGeotypeCV *string `json:"samplingFeatureGeoTypeCV"`
	// nullable
	FeatureGeometryWKT *string `json:"featureGeometryWKT"`
	// nullable
	ElevationM *float64 `json:"elevationM"`
	// nullable
	ElevationDatumCV *string `json:"elevationDatumCV"`
	// nullable
	ElevationPrecision *float64 `json:"elevationPrecision"`
	// nullable
	ElevationPrecisionComment *string `json:"elevationPrecisionComment"`
}

type SampleResponse struct {
	NumItems int      `json:"numItems"`
	Data     []Sample `json:"data"`
}

type SampleByFilters struct {
	SampleID          int       `json:"samplingfeatureid"`
	SampleName        string    `json:"sampleName"`
	Batches           []*int    `json:"batches"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	PublicationYear   *int      `json:"publicationYear"`
	DOI               *string   `json:"doi"`
	Authors           []*Author `json:"authors"`
	Minerals          []*string `json:"minerals"`
	HostMinerals      []*string `json:"hostMinerals"`
	InclusionMinerals []*string `json:"inclusionMinerals"`
	RockClasses       []*string `json:"rockClasses"`
	RockTypes         []*string `json:"rockTypes"`
	InclusionTypes    []*string `json:"inclusionTypes"`
	GeologicalSetting []*string `json:"geologicalSettings"`
	GeologicalAge     []*string `json:"geologicalAges"`
	GeologicalAgesMin []*string `json:"geologicalAgesMin"`
	GeologicalAgesMax []*string `json:"geologicalAgesMax"`
	TotalCount        int       `json:"totalCount"`
}

type SampleByFiltersData struct {
	SampleID   int     `json:"sampleID"`
	SampleName string  `json:"sampleName"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	// nullable
	PublicationYear *int `json:"publicationYear"`
	// nullable
	DOI     *string   `json:"doi"`
	Authors []*Author `json:"authors"`
	// only filled for material="Mineral"
	Minerals []*string `json:"minerals"`
	// only filled for material="Inclusion"
	HostMinerals []*string `json:"hostMinerals"`
	// only filled for material="Inclusion"
	InclusionMinerals []*string `json:"inclusionMinerals"`
	// only filled for material="WholeRock or Glass"
	RockClasses []*string `json:"rockClasses"`
	RockTypes   []*string `json:"rockTypes"`
	// only filled for material="Inclusion"
	InclusionTypes    []*string `json:"inclusionTypes"`
	GeologicalSetting []*string `json:"geologicalSettings"`
	GeologicalAge     []*string `json:"geologicalAges"`
	GeologicalAgesMin []*string `json:"geologicalAgesMin"`
	GeologicalAgesMax []*string `json:"geologicalAgesMax"`
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
