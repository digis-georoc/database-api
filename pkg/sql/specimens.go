package sql

const GetSpecimenTypesQuery = `
select distinct s.specimentypecv as specimentype
from odm2.specimens s
`
