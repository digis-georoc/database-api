package model

type FullData struct {
	SampleID           int        `json:"sampleID"`
	UniqueID           *string    `json:"uniqueID"`
	References         []Citation `json:"references"`
	SampleName         *string    `json:"sampleName"`
	LocationNames      []*string  `json:"locationNames"`
	LocationTypes      []*string  `json:"locationTypes"`
	ElevationMin       *string    `json:"elevationMin"`
	ElevationMax       *string    `json:"elevationMax"`
	LandOrSea          *string    `json:"landOrSea"`
	RockType           *string    `json:"rockType"`
	RockClass          *string    `json:"rockClass"`
	RockTextures       []*string  `json:"rockTextures"`
	AgeMin             *int       `json:"ageMin"`
	AgeMax             *int       `json:"ageMax"`
	EruptionDate       *string    `json:"eruptionDate"`
	GeologicalAge      *string    `json:"geologicalAge"`
	Minerals           []*string  `json:"minerals"`
	HostMinerals       []*string  `json:"hostMinerals"`
	LocationNum        *int       `json:"locationNum"`
	Latitude           *float32   `json:"latitude"`
	Longitude          *float32   `json:"longitude"`
	LatitudeMin        *string    `json:"latitudeMin"`
	LatitudeMax        *string    `json:"latitudeMax"`
	LongitudeMin       *string    `json:"longitudeMin"`
	LongitudeMax       *string    `json:"longitudeMax"`
	TectonicSetting    *string    `json:"tectonicSetting"`
	Method             []*string  `json:"method"`
	Comment            []*string  `json:"comment"`
	Institutions       []*string  `json:"institutions"`
	Results            []*Result  `json:"results"`
	Alterations        []*string  `json:"alterations"`
	SamplingTechniques []*string  `json:"samplingTechniques"`
	DrillDepthMin      []*string  `json:"drillDepthMin"`
	DrillDepthMax      []*string  `json:"drillDepthMax"`
	BatchData          []*Batch   `json:"batchData"`
}

type Batch struct {
	BatchID            *int      `json:"batchID"`
	BatchName          *string   `json:"batchName"`
	SampleID           *int      `json:"sampleID"`
	Crystal            *string   `json:"crystal"`
	SpecimenMedium     *string   `json:"specimenMedium"`
	RockTypes          []*string `json:"rockTypes"`
	RockClasses        []*string `json:"rockClasses"`
	Minerals           []*string `json:"minerals"`
	Materials          []*string `json:"materials"`
	InclusionTypes     []*string `json:"inclusionTypes"`
	RimOrCoreInclusion []*string `json:"rimOrCoreInclusion"`
	RimOrCoreMineral   []*string `json:"rimOrCoreMineral"`
}

type FullDataResponse struct {
	NumItems int        `json:"numItems"`
	Data     []FullData `json:"data"`
}
