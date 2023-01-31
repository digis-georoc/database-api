package sql

const RockClassQuery = `
select distinct t.taxonomicclassifiername as name 
from odm2.taxonomicclassifiers t
where t.taxonomicclassifiertypecv = 'Lithology'
order by t.taxonomicclassifiername 
`

const RockTypeQuery = `
select distinct t.taxonomicclassifiername as name
from odm2.taxonomicclassifiers t
where t.taxonomicclassifiertypecv = 'Rock'
order by t.taxonomicclassifiername 
`

const MineralQuery = `
select distinct t.taxonomicclassifiercommonname as name
from odm2.taxonomicclassifiers t
where t.taxonomicclassifierdescription = 'Mineral Classification from GEOROC'
order by t.taxonomicclassifiercommonname
`
