package model

type Statistics struct {
	NumCitations int `json:"numCitations"`
	NumSamples   int `json:"numSamples"`
	NumAnalyses  int `json:"numAnalyses"`
	NumResults   int `json:"numResults"`
}
