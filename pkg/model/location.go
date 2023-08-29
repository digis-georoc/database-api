package model

type Location struct {
	Name string `json:"name"`
}

type LocationResponse struct {
	NumItems int        `json:"numItems"`
	Data     []Location `json:"data"`
}
