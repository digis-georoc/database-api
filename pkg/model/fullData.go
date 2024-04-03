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
	SampleName *string `json:"sampleName"`
	// nullable
	LocationNames []*string `json:"locationNames"`
	// nullable
	LocationTypes []*string `json:"locationTypes"`
	// nullable
	LocationComments []*string `json:"locationComments"`
	// nullable
	ElevationMin *string `json:"elevationMin"`
	// nullable
	ElevationMax *string `json:"elevationMax"`
	// nullable
	LandOrSea *string `json:"landOrSea"`
	// nullable
	RockTypes []*FullDataTaxonomicClassifier `json:"rockTypes"`
	// nullable
	RockClasses []*FullDataTaxonomicClassifier `json:"rockClasses"`
	// nullable
	RockTextures []*string `json:"rockTextures"`
	// nullable
	AgeMin *int `json:"ageMin"`
	// nullable
	AgeMax *int `json:"ageMax"`
	// nullable
	EruptionDate *string `json:"eruptionDate"`
	// nullable
	GeologicalAge *string `json:"geologicalAge"`
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
	TectonicSetting *string `json:"tectonicSetting"`
	// nullable
	Methods []*string `json:"methods"`
	// nullable
	Comments []*string `json:"comments"`
	// nullable
	Institutions []*string `json:"institutions"`
	// nullable
	Results []*Result `json:"results"`
	// nullable
	Alteration *string `json:"alteration"`
	// nullable
	SamplingTechnique *string `json:"samplingTechnique"`
	// nullable
	DrillDepthMin *string `json:"drillDepthMin"`
	// nullable
	DrillDepthMax *string `json:"drillDepthMax"`
	// nullable
	BatchData []*Batch `json:"batchData"`
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
	SpecimenMedium *string `json:"specimenMedium"`
	// nullable
	Minerals []*FullDataTaxonomicClassifier `json:"minerals"`
	// nullable
	HostMinerals []*FullDataTaxonomicClassifier `json:"hostMinerals"`
	// nullable
	InclusionMinerals []*FullDataTaxonomicClassifier `json:"inclusionMinerals"`
	// nullable
	Materials []*string `json:"materials"`
	// nullable
	InclusionTypes []*string `json:"inclusionTypes"`
	// nullable
	RimOrCoreInclusion *string `json:"rimOrCoreInclusion"`
	// nullable
	RimOrCoreMineral *string `json:"rimOrCoreMineral"`
	// nullable
	Results []*Result `json:"results"`
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
