// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

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
