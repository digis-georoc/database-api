package sql

const GetSpecimenTypesQuery = `
select distinct s.specimentypecv as specimentypecv
from odm2.specimens s
`
