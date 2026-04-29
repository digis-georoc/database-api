// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package model

import (
	"fmt"
	"regexp"
)

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

// ParseChemQuery takes a chemistry query DSL string and parses it into a ChemQuery structure
func ParseChemQuery(query string) (ChemQuery, error) {
	chemQuery := ChemQuery{}
	expressionRegex := regexp.MustCompile(`\(([\w]+)?,([\w\d]+)?,([\d\.]+)?,([\d\.]+)?\)`)
	matches := expressionRegex.FindAllStringSubmatch(query, -1)
	if len(matches) == 0 {
		return chemQuery, fmt.Errorf("can not parse chemical query")
	}
	for i, match := range matches {
		junctor := CQ_JUNCTOR_AND
		if i == 0 {
			// first expression gets no junctor
			junctor = CQ_JUNCTOR_NONE
		}
		expr := CQExpression{
			Junctor:  junctor,
			Type:     match[1],
			Element:  match[2],
			MinValue: match[3],
			MaxValue: match[4],
		}
		// omit expressions without type or element as they make no sense
		if expr.Type == "" && expr.Element == "" {
			continue
		}
		chemQuery.Expressions = append(chemQuery.Expressions, expr)
	}
	return chemQuery, nil
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
