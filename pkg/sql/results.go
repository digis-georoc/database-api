package sql

const ElementsQuery = `
select distinct v.variablecode as name
from odm2.variables v 
`

const ElementTypesQuery = `
select distinct v.variabletypecode as name
from odm2.variables v 
`
