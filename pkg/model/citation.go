package model

type Citation struct {
	CitationID      int      `json:"citationID"`
	Title           string   `json:"title"`
	Publisher       string   `json:"publisher"`
	Publicationyear int      `json:"publicationYear"`
	CitationLink    string   `json:"citationLink"`
	Journal         string   `json:"journal"`
	Volume          string   `json:"volume"`
	Issue           string   `json:"issue"`
	FirstPage       string   `json:"firstPage"`
	LastPage        string   `json:"lastPage"`
	BookTitle       string   `json:"bookTitle"`
	Editors         string   `json:"editors"`
	Authors         []Person `json:"authors"`
}

type CitationResponse struct {
	NumItems int        `json:"numItems"`
	Data     []Citation `json:"data"`
}
