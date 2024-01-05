package sql

const GetSpecimenTypesQuery = `
select distinct s.specimentypecv as specimentypecv
from odm2.specimens s
`

const GetRandomSpecimensQuery = `
with params as (
select min(s.samplingfeatureid) as min_id, max(s.samplingfeatureid) as max_id
from odm2.specimens s 
)
select samplingFeatureID,
specimenTypeCV,
specimenMediumCV,
isFieldSpecimen
from (
select trunc(p.min_id+random() * p.max_id) as id from
params p,
generate_series(1, $1 + ($1 / 10))
) rand
join odm2.specimens s on s.samplingfeatureid = rand.id
limit $1
`

const GetGeoAgesQuery = `
select distinct specimengeolage as name
from odm2.specimenages
`

const GetGeoAgePrefixesQuery = `
select distinct specimengeolageprefix as name
from odm2.specimenages
`
