// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

type Citation struct {
	CitationID int `json:"citationID"`
	// nullable
	Title *string `json:"title"`
	// nullable
	Publisher *string `json:"publisher"`
	// nullable
	Publicationyear *int `json:"publicationYear"`
	// nullable
	CitationLink *string `json:"citationLink"`
	// nullable
	Journal *string `json:"journal"`
	// nullable
	Volume *string `json:"volume"`
	// nullable
	Issue *string `json:"issue"`
	// nullable
	FirstPage *string `json:"firstPage"`
	// nullable
	LastPage *string `json:"lastPage"`
	// nullable
	BookTitle *string `json:"bookTitle"`
	// nullable
	Editors *string  `json:"editors"`
	Authors []Author `json:"authors"`
	// nullable
	DOI *string `json:"doi"`
}

type Author struct {
	PersonID int `json:"personID"`
	// nullable
	PersonFirstName *string `json:"personFirstName"`
	// nullable
	PersonLastName *string `json:"personLastName"`
	// nullable
	AuthorOrder *int `json:"authorOrder"`
}

type CitationResponse struct {
	NumItems int        `json:"numItems"`
	Data     []Citation `json:"data"`
}
