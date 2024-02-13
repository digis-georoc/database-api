// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

type Citation struct {
	CitationID      int      `json:"citationID"`
	Title           *string  `json:"title"`
	Publisher       *string  `json:"publisher"`
	Publicationyear *int     `json:"publicationYear"`
	CitationLink    *string  `json:"citationLink"`
	Journal         *string  `json:"journal"`
	Volume          *string  `json:"volume"`
	Issue           *string  `json:"issue"`
	FirstPage       *string  `json:"firstPage"`
	LastPage        *string  `json:"lastPage"`
	BookTitle       *string  `json:"bookTitle"`
	Editors         *string  `json:"editors"`
	Authors         []Author `json:"authors"`
	DOI             *string  `json:"doi"`
}

type Author struct {
	PersonID        int     `json:"personID"`
	PersonFirstName *string `json:"personFirstName"`
	PersonLastName  *string `json:"personLastName"`
	AuthorOrder     *int    `json:"authorOrder"`
}

type CitationResponse struct {
	NumItems int        `json:"numItems"`
	Data     []Citation `json:"data"`
}
