// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

const (
	CQ_JUNCTOR_NONE = "%"
	CQ_JUNCTOR_AND  = "^"
	CQ_JUNCTOR_OR   = "|"
)

// cqSQLMap maps chemical query symbols to sql
var CQSQLMap = map[string]string{
	CQ_JUNCTOR_NONE: "from (",
	CQ_JUNCTOR_AND:  "join (",
	CQ_JUNCTOR_OR:   "full join (",
}

type Element struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Unit  string `json:"unit"`
}

type ElementResponse struct {
	NumItems int       `json:"numItems"`
	Data     []Element `json:"data"`
}

type ElementType struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ElementTypeResponse struct {
	NumItems int           `json:"numItems"`
	Data     []ElementType `json:"data"`
}

type CQExpression struct {
	Junctor  string
	Type     string
	Element  string
	MinValue string
	MaxValue string
}

type ChemQuery struct {
	Expressions []CQExpression
}

type Result struct {
	// nullable
	ItemName *string `json:"itemName"`
	// nullable
	ItemGroup *string `json:"itemGroup"`
	// nullable
	Medium *string `json:"medium"`
	// nullable
	ValueCount *int       `json:"valueCount"`
	Standards  []Standard `json:"standards"`
	// nullable
	Value *float64 `json:"value"`
	// nullable
	Unit *string `json:"unit"`
	// nullable
	Method *string `json:"method"`
}

type Standard struct {
	StandardName     string  `json:"standardName"`
	StandardValue    float64 `json:"standardValue"`
	StandardVariable string  `json:"standardVariable"`
	StandardUnit     string  `json:"standardUnit"`
}

type Measurement struct {
	Element string  `json:"element"`
	Value   float64 `json:"value"`
	Unit    string  `json:"unit"`
}
