// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

type FullData struct {
	SampleID int `json:"sampleID"`
	// nullable
	UniqueID   *string    `json:"uniqueID"`
	References []Citation `json:"references"`
	// nullable
	SampleName       *string   `json:"sampleName"`
	LocationNames    []*string `json:"locationNames"`
	LocationTypes    []*string `json:"locationTypes"`
	LocationComments []*string `json:"locationComments"`
	// nullable
	ElevationMin *string `json:"elevationMin"`
	// nullable
	ElevationMax *string `json:"elevationMax"`
	// nullable
	LandOrSea    *string                        `json:"landOrSea"`
	RockTypes    []*FullDataTaxonomicClassifier `json:"rockTypes"`
	RockClasses  []*FullDataTaxonomicClassifier `json:"rockClasses"`
	RockTextures []*string                      `json:"rockTextures"`
	// nullable
	AgeMin *float64 `json:"ageMin"`
	// nullable
	AgeMax *float64 `json:"ageMax"`
	// nullable
	EruptionDate *string `json:"eruptionDate"`
	// nullable
	GeologicalAge *string `json:"geologicalAge"`
	// nullable
	GeologicalAgePrefix *string `json:"geologicalAgePrefix"`
	// nullable
	LocationNum *int `json:"locationNum"`
	// nullable
	Latitude *float32 `json:"latitude"`
	// nullable
	Longitude *float32 `json:"longitude"`
	// nullable
	LatitudeMin *string `json:"latitudeMin"`
	// nullable
	LatitudeMax *string `json:"latitudeMax"`
	// nullable
	LongitudeMin *string `json:"longitudeMin"`
	// nullable
	LongitudeMax *string `json:"longitudeMax"`
	// nullable
	TectonicSetting *string   `json:"tectonicSetting"`
	Methods         []*string `json:"methods"`
	Comments        []*string `json:"comments"`
	Institutions    []*string `json:"institutions"`
	Results         []*Result `json:"results"`
	// nullable
	Alteration *string `json:"alteration"`
	// nullable
	AlterationType *string `json:"alterationType"`
	// nullable
	SamplingTechnique *string `json:"samplingTechnique"`
	// nullable
	DrillDepthMin *string `json:"drillDepthMin"`
	// nullable
	DrillDepthMax *string  `json:"drillDepthMax"`
	BatchData     []*Batch `json:"batchData"`
}

type Batch struct {
	// nullable
	BatchID *int `json:"batchID"`
	// nullable
	BatchName *string `json:"batchName"`
	// nullable
	SampleID *int `json:"sampleID"`
	// nullable
	Crystal *string `json:"crystal"`
	// nullable
	SpecimenMedium    *string                        `json:"specimenMedium"`
	Minerals          []*FullDataTaxonomicClassifier `json:"minerals"`
	HostMinerals      []*FullDataTaxonomicClassifier `json:"hostMinerals"`
	InclusionMinerals []*FullDataTaxonomicClassifier `json:"inclusionMinerals"`
	// nullable
	Material       *string   `json:"material"`
	InclusionTypes []*string `json:"inclusionTypes"`
	// nullable
	RimOrCoreInclusion *string `json:"rimOrCoreInclusion"`
	// nullable
	RimOrCoreMineral *string   `json:"rimOrCoreMineral"`
	Results          []*Result `json:"results"`
	// nullable
	TASData *DiagramData `json:"tasData" db:"-"`
}

type DiagramData struct {
	XAxisLabel string      `json:"xAxisLabel"`
	YAxisLabel string      `json:"yAxisLabel"`
	Values     [][]float64 `json:"values"`
}

type FullDataResponse struct {
	NumItems int        `json:"numItems"`
	Data     []FullData `json:"data"`
}

type FullDataTaxonomicClassifier struct {
	Value string  `json:"value"`
	Label *string `json:"label"`
	ID    int     `json:"id"`
}
