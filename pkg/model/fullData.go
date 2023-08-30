package model

type FullData struct {
	SampleNum       int           `json:"sampleNum"`
	UniqueID        string        `json:"uniqueID"`
	Batches         []*int        `json:"batches"`
	References      []interface{} `json:"references"`
	SampleName      string        `json:"sampleName"`
	LocationNames   []string      `json:"locationNames"`
	LocationTypes   []string      `json:"locationTypes"`
	ElevationMin    string        `json:"elevationMin"`
	ElevationMax    string        `json:"elevationMax"`
	LandOrSea       string        `json:"landOrSea"`
	RockTypes       []string      `json:"rockTypes"`
	RockClasses     []string      `json:"rockClasses"`
	RockTextures    []string      `json:"rockTextures"`
	AgeMin          *int          `json:"ageMin"`
	AgeMax          *int          `json:"ageMax"`
	Materials       []string      `json:"materials"`
	Minerals        []string      `json:"minerals"`
	InclusionTypes  []string      `json:"inclusionTypes"`
	LocationNum     *int          `json:"locationNum"`
	Latitude        *float32      `json:"latitude"`
	Longitude       *float32      `json:"longitude"`
	LatitudeMin     string        `json:"latitudeMin"`
	LatitudeMax     string        `json:"latitudeMax"`
	LongitudeMin    string        `json:"longitudeMin"`
	LongitudeMax    string        `json:"longitudeMax"`
	TectonicSetting string        `json:"tectonicSetting"`
	Method          []string      `json:"method"`
	Comment         []string      `json:"comment"`
	Institutions    []string      `json:"institutions"`
	ItemName        []string      `json:"itemName"`
	ItemGroup       []string      `json:"itemGroup"`
	StandardNames   []string      `json:"standardNames"`
	StandardValues  []*float32    `json:"standardValues"`
	Values          []*float32    `json:"values"`
	Units           []string      `json:"units"`
}

type FullDataResponse struct {
	NumItems int        `json:"numItems"`
	Data     []FullData `json:"data"`
}
