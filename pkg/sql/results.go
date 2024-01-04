package sql

const ElementsQuery = `
select v.variablecode as value,
v.variablecode as label,
'' as unit
from odm2.variables v
`

const ElementTypesQuery = `
select distinct v.variabletypecode as value,
v.variabletypecv as label
from odm2.variables v 
where v.variabletypecv !='Rock mode'
`
