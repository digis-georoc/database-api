package model

type TaxonomicClassifier struct {
	Value string  `json:"value"`
	Label *string `json:"label"`
	Count int     `json:"count"`
}

type TaxonomicClassifierResponse struct {
	NumItems int                   `json:"numItems"`
	Data     []TaxonomicClassifier `json:"data"`
}
