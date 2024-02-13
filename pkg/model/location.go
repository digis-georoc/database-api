// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

type Location struct {
	Name string `json:"name"`
}

type LocationResponse struct {
	NumItems int        `json:"numItems"`
	Data     []Location `json:"data"`
}
