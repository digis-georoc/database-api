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
	Name string `json:"name"`
}

type ElementResponse struct {
	NumItems int       `json:"numItems"`
	Data     []Element `json:"data"`
}

type ElementType struct {
	Name string `json:"name"`
}

type ElementTypeResponse struct {
	NumItems int           `json:"numItems"`
	Data     []ElementType `json:"data"`
}

type CQExpression struct {
	Junctor  string
	Element  string
	Operator string
	Value    string
}

type ChemQuery struct {
	Expressions []CQExpression
}
