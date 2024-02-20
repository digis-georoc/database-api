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
	ElevationMin *string `json:"elevationMin"`
	// nullable
	ElevationMax *string `json:"elevationMax"`
	// nullable
	LandOrSea *string `json:"landOrSea"`
	// nullable
	RockType *string `json:"rockType"`
	// nullable
	RockClass *string `json:"rockClass"`
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
	Minerals []*string `json:"minerals"`
	// nullable
	HostMinerals []*string `json:"hostMinerals"`
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
	Method []*string `json:"method"`
	// nullable
	Comment []*string `json:"comment"`
	// nullable
	Institutions []*string `json:"institutions"`
	// nullable
	Results []*Result `json:"results"`
	// nullable
	Alterations []*string `json:"alterations"`
	// nullable
	SamplingTechniques []*string `json:"samplingTechniques"`
	// nullable
	DrillDepthMin []*string `json:"drillDepthMin"`
	// nullable
	DrillDepthMax []*string `json:"drillDepthMax"`
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
	RockTypes []*string `json:"rockTypes"`
	// nullable
	RockClasses []*string `json:"rockClasses"`
	// nullable
	Minerals []*string `json:"minerals"`
	// nullable
	Materials []*string `json:"materials"`
	// nullable
	InclusionTypes []*string `json:"inclusionTypes"`
	// nullable
	RimOrCoreInclusion *string `json:"rimOrCoreInclusion"`
	// nullable
	RimOrCoreMineral *string `json:"rimOrCoreMineral"`
}

type FullDataResponse struct {
	NumItems int        `json:"numItems"`
	Data     []FullData `json:"data"`
}
