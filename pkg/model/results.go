package model

type Element struct {
	Name string `json:"name"`
}

type ElementResponse struct {
	NumItems int       `json:"numItems"`
	Data     []Element `json:"data"`
}

type ElementType struct {
	Name string `json:"name"`
}

type ElementTypeResponse struct {
	NumItems int           `json:"numItems"`
	Data     []ElementType `json:"data"`
}
