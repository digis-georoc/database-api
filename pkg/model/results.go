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
	ItemName      *string  `json:"itemName"`
	ItemGroup     *string  `json:"itemGroup"`
	StandardName  *string  `json:"standardName"`
	StandardValue *float32 `json:"standardValue"`
	Value         *float32 `json:"value"`
	Unit          *string  `json:"unit"`
}
