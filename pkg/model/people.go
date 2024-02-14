// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

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
