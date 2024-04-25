// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

type Statistics struct {
	NumCitations int    `json:"numCitations"`
	NumSamples   int    `json:"numSamples"`
	NumAnalyses  int    `json:"numAnalyses"`
	NumResults   int    `json:"numResults"`
	LatestDate   string `json:"latestDate"`
}
