package model

type Person struct {
	PersonID        int     `json:"personID"`
	PersonFirstName *string `json:"personFirstName"`
	PersonLastName  *string `json:"personLastName"`
}

type PeopleResponse struct {
	NumItems int      `json:"numItems"`
	Data     []Person `json:"data"`
}
