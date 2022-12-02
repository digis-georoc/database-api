package model

type Citation struct {
	CitationID      int
	Title           string
	Publisher       string
	Publicationyear int
	CitationLink    string
	Journal         string
	Volume          string
	Issue           string
	FirstPage       string
	LastPage        string
	BookTitle       string
	Editors         string
	Authors         []People
}
