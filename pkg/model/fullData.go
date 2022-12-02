package model

type FullData struct {
	Sample_Num       int
	UniqueId         string `json:"unique_id"`
	Batches          []int
	References       interface{}
	SampleIds        []string
	Location_Names   [][]string
	Location_Types   [][]string
	Loc_Data         interface{}
	Elevation_Min    int
	Elevation_Max    int
	Land_Or_Sea      [][]string
	Rock_Type        string
	Rock_Class       string
	Rock_Texture     string
	Age_Min          int
	Age_Max          int
	Material         []string
	Mineral          []string
	Inclusion_Type   []string
	Location_Num     int
	Latitude         float32
	Longitude        float32
	Latitude_Min     string
	Latitude_Max     string
	Longitude_Min    string
	Longitude_Max    string
	Tectonic_Setting [][]string
	Method           []string
	Comment          []string
	Institution      [][]string
	Item_Name        [][]string
	Item_Group       [][]string
	Standard_Names   [][]string
	Standard_Values  [][]float32
	Values           [][]float32
	Units            [][]string
}
